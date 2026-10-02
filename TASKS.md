# Tasks

One-line tasks a fleet Bot can take with almost no context. Each one is done when
`just check` passes with a new unit test that proves the behaviour.

Tasks 1-4 are independent: put the new handler, any new `Store` method and its
`init()` that calls `register(...)` in a new `<feature>.go`, and the tests in a new
`<feature>_test.go`. They then never edit the same lines, so they land in any order
without conflicts.

Tasks 5 and 6 deliberately conflict: both add validation at the top of `handlePut`
in `server.go`, so whichever lands second has to rebase onto the first.

1. **DELETE a key.** `DELETE /kv/{key}` removes the key and returns 204; a later `GET` returns 404, and deleting a missing key returns 404.
2. **Health endpoint.** `GET /healthz` returns 200 with body `ok`.
3. **Count keys.** `GET /count` returns 200 with the number of stored keys as a plain decimal body, e.g. `3`.
4. **List keys.** `GET /kv` returns 200 with a JSON array of every key, sorted, e.g. `["a","b"]`; an empty store returns `[]`.
5. **Limit value size.** In `handlePut`, a body over 1024 bytes returns 413 and stores nothing.
6. **Validate keys.** In `handlePut`, a key longer than 64 bytes returns 400 and stores nothing.
