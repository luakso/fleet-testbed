// Command fleet-testbed serves a tiny in-memory key-value store over HTTP.
//
//	GET /kv/{key}   returns the value, or 404
//	PUT /kv/{key}   stores the request body, returns 204
//
// It only listens on a loopback address and stops gracefully on SIGINT or SIGTERM.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "loopback address to listen on; port 0 picks a free port")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, *addr); err != nil {
		log.Fatal(err)
	}
}

// run serves until ctx is done, then shuts down gracefully. Once listening it
// prints "listening on http://ADDR" to stdout so callers can find a port-0 address.
func run(ctx context.Context, addr string) error {
	if err := checkLoopback(addr); err != nil {
		return err
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := &http.Server{Handler: NewHandler(NewStore()), ReadHeaderTimeout: 5 * time.Second}
	fmt.Printf("listening on http://%s\n", ln.Addr())

	errc := make(chan error, 1)
	go func() { errc <- srv.Serve(ln) }()
	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	if err := <-errc; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// checkLoopback rejects any address whose host is not a loopback IP or localhost.
func checkLoopback(addr string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	if host == "localhost" {
		return nil
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		return nil
	}
	return fmt.Errorf("refusing non-loopback address %q", addr)
}
