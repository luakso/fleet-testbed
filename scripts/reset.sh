#!/usr/bin/env bash
# Return this clone to the `start` tag so every fleet run begins from the same code.
#
# Refuses unless run inside a git clone whose origin is luakso/fleet-testbed. It only
# changes that clone: it force-switches the local `main` branch to `start` and removes
# untracked and ignored files inside the clone. It never touches GitHub's main.
set -euo pipefail

top=$(git rev-parse --show-toplevel 2>/dev/null) || {
	echo "reset: not inside a git clone; refusing" >&2
	exit 1
}
script_top=$(cd "$(dirname "$0")" && git rev-parse --show-toplevel)
if [ "$top" != "$script_top" ]; then
	echo "reset: run this from the clone it belongs to ($script_top), not $top; refusing" >&2
	exit 1
fi

origin=$(git -C "$top" remote get-url origin 2>/dev/null) || {
	echo "reset: clone has no origin remote; refusing" >&2
	exit 1
}
case "$origin" in
https://github.com/luakso/fleet-testbed | https://github.com/luakso/fleet-testbed.git | \
	git@github.com:luakso/fleet-testbed | git@github.com:luakso/fleet-testbed.git | \
	ssh://git@github.com/luakso/fleet-testbed | ssh://git@github.com/luakso/fleet-testbed.git) ;;
*)
	echo "reset: origin is $origin, not luakso/fleet-testbed; refusing" >&2
	exit 1
	;;
esac

cd "$top"
if ! git fetch --quiet --force origin "+refs/tags/start:refs/tags/start"; then
	echo "reset: could not fetch the start tag from origin; using the local copy if present" >&2
fi
git rev-parse --verify --quiet "refs/tags/start^{commit}" >/dev/null || {
	echo "reset: no start tag found; refusing" >&2
	exit 1
}

git switch --quiet --force --create main start
git clean --quiet --force -d -x
echo "reset: $top is at start ($(git rev-parse --short HEAD))"
