// Command smoke builds the server, starts it on a free loopback port, makes real
// requests against it and stops it with SIGTERM. It removes the binary it built
// and never leaves the server running, even when a step fails.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "smoke: FAIL:", err)
		os.Exit(1)
	}
	fmt.Println("smoke: ok")
}

func run() error {
	dir, err := os.MkdirTemp("", "fleet-testbed-smoke-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	bin := dir + "/server"
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Stdout, build.Stderr = os.Stderr, os.Stderr
	if err := build.Run(); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	srv := exec.Command(bin, "-addr", "127.0.0.1:0")
	srv.Stderr = os.Stderr
	out, err := srv.StdoutPipe()
	if err != nil {
		return err
	}
	if err := srv.Start(); err != nil {
		return err
	}
	exited := make(chan error, 1)
	go func() { exited <- srv.Wait() }()
	stopped := false
	defer func() {
		// The server is still running here only when a step failed.
		if !stopped {
			srv.Process.Kill()
			<-exited
		}
	}()

	base, err := readAddr(out)
	if err != nil {
		return err
	}
	if err := exercise(base); err != nil {
		return err
	}

	if err := srv.Process.Signal(syscall.SIGTERM); err != nil {
		return err
	}
	select {
	case err := <-exited:
		stopped = true
		if err != nil {
			return fmt.Errorf("server did not stop cleanly: %w", err)
		}
	case <-time.After(10 * time.Second):
		return errors.New("server did not stop within 10s of SIGTERM")
	}
	return nil
}

// readAddr waits for the server's "listening on URL" line and returns URL.
func readAddr(r io.Reader) (string, error) {
	line := make(chan string, 1)
	go func() {
		s := bufio.NewScanner(r)
		if s.Scan() {
			line <- s.Text()
		}
		close(line)
		io.Copy(io.Discard, r)
	}()
	select {
	case l, ok := <-line:
		url, found := strings.CutPrefix(l, "listening on ")
		if !ok || !found {
			return "", fmt.Errorf("unexpected server output %q", l)
		}
		return url, nil
	case <-time.After(10 * time.Second):
		return "", errors.New("server did not report its address within 10s")
	}
}

// exercise stores a value, reads it back and checks a missing key is a 404.
func exercise(base string) error {
	c := &http.Client{Timeout: 5 * time.Second}

	req, _ := http.NewRequest(http.MethodPut, base+"/kv/smoke", strings.NewReader("ok"))
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("PUT: got %s, want 204", resp.Status)
	}

	resp, err = c.Get(base + "/kv/smoke")
	if err != nil {
		return err
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK || string(body) != "ok" {
		return fmt.Errorf("GET: got %s %q, want 200 %q", resp.Status, body, "ok")
	}

	resp, err = c.Get(base + "/kv/missing")
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("GET missing: got %s, want 404", resp.Status)
	}
	return nil
}
