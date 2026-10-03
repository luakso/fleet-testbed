#!/usr/bin/env bash
# Check that reset.sh returns a throwaway clone to `start`, both when no local `main`
# exists and when one does. Offline: https fetches are disabled, so reset.sh falls back
# to the local start tag and never reaches GitHub.
set -euo pipefail

here=$(cd "$(dirname "$0")" && pwd)
tmp=$(mktemp -d)
trap 'rm -rf -- "$tmp"' EXIT
export GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=protocol.https.allow GIT_CONFIG_VALUE_0=never
export GIT_AUTHOR_NAME=test GIT_AUTHOR_EMAIL=test@example.com
export GIT_COMMITTER_NAME=test GIT_COMMITTER_EMAIL=test@example.com

clone="$tmp/clone"
git init --quiet --initial-branch=work "$clone"
mkdir "$clone/scripts"
cp "$here/reset.sh" "$clone/scripts/reset.sh"
git -C "$clone" add scripts
git -C "$clone" commit --quiet -m start
git -C "$clone" tag start
git -C "$clone" remote add origin https://github.com/luakso/fleet-testbed.git
start=$(git -C "$clone" rev-parse start)

expect_start() {
	[ "$(git -C "$clone" symbolic-ref --short HEAD)" = main ] || { echo "reset_test: $1: not on main" >&2; exit 1; }
	[ "$(git -C "$clone" rev-parse HEAD)" = "$start" ] || { echo "reset_test: $1: main is not at start" >&2; exit 1; }
	[ -z "$(git -C "$clone" status --porcelain --ignored)" ] || { echo "reset_test: $1: clone not clean" >&2; exit 1; }
}

# No local main yet.
echo change >"$clone/later.txt"
git -C "$clone" add later.txt
git -C "$clone" commit --quiet -m later
(cd "$clone" && ./scripts/reset.sh >/dev/null)
expect_start "without main"

# main exists and has moved past start, with stray files.
echo change >"$clone/later.txt"
git -C "$clone" add later.txt
git -C "$clone" commit --quiet -m later
echo stray >"$clone/stray.txt"
(cd "$clone" && ./scripts/reset.sh >/dev/null)
expect_start "with main"

echo "reset_test: ok"
