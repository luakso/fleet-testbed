# fleet-testbed

Small, resettable test project for exercising [fleet](https://github.com/luakso/fleet)
end to end. It is a dependency-free Go key-value HTTP server (standard library only),
small enough that a Bot can take a task with almost no context, and resettable so
every run starts from the same code.

## The server

```sh
go run . -addr 127.0.0.1:8080   # loopback addresses only; port 0 picks a free port
curl -X PUT --data hello http://127.0.0.1:8080/kv/greeting   # 204
curl http://127.0.0.1:8080/kv/greeting                       # hello
```

It stops gracefully on SIGINT or SIGTERM.

## Checking

```sh
just check
```

This checks formatting, vets, builds, runs the unit tests, then builds the server into a
temporary directory, starts it on a free loopback port, makes real requests and stops it.
It leaves no process or file behind, and needs nothing beyond loopback networking, so it
runs inside Claude Code's sandbox. GitHub Actions runs the same command on pushes and pull
requests to `main`.

## Tasks

[TASKS.md](TASKS.md) lists six one-line tasks with what "done" means. Tasks 5 and 6
deliberately edit the same code, so landing one forces the other to rebase.

## Sandbox-widening fixture

[fixtures/claude-widening-settings.json](fixtures/claude-widening-settings.json) is a
Claude Code settings file that tries to widen the sandbox: it allows reading `~/.ssh`,
writing the whole home directory, all Unix sockets, unsandboxed and excluded commands,
and broad allow rules. It is kept as a fixture, not at `.claude/settings.json`, so opening
Claude Code in this repository never loads it. A test harness copies it into a throwaway
clone's `.claude/settings.json` to check that fleet does not honour it. `.claude/` is
ignored so the copy is never committed.

## Resetting

The starting point is the `start` tag on `main`.

```sh
just reset
```

returns a local clone to `start`: it fetches the tag, force-switches the local `main`
branch to it and removes untracked and ignored files (including a copied
`.claude/settings.json`). It refuses unless it runs inside a clone whose `origin` is
`luakso/fleet-testbed`, and it never touches anything outside that clone or on GitHub.
Other local branches are left alone.

Resetting GitHub's `main` back to `start` is a manual step for the captain only:

```sh
git push --force origin start:main
```
