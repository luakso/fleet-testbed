# List recipes.
default:
    @just --list

# Format check, vet, build, unit tests, then a real start/request/stop smoke run.
check:
    test -z "$(gofmt -l .)" || { gofmt -l .; echo "gofmt needed on the files above"; exit 1; }
    go vet ./...
    go build -o /dev/null .
    go test ./...
    go run ./internal/smoke
    ./scripts/reset_test.sh

# Return this clone to the `start` tag. Refuses outside a luakso/fleet-testbed clone.
reset:
    ./scripts/reset.sh
