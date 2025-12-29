# Litestream for the `ncruces/go-sqlite3` driver

This is a fork focusing on **library** use, which uses the
[`ncruces/go-sqlite3`](https://github.com/ncruces/go-sqlite3) driver.

If you _can_ use the `litestream` CLI for streaming replication, please do.
The `ncruces/go-sqlite3` driver is [compatible](https://github.com/ncruces/go-sqlite3/discussions/68)
with the CLI on all major platforms.

Otherwise, check the [examples](_examples/library/) for how to use.

## Lightweight read-replicas

This also implements the `"litestream"` SQLite VFS that offers
Litestream [lightweight read-replicas](https://fly.io/blog/litestream-revamped/#lightweight-read-replicas)
for the `ncruces/go-sqlite3` driver.

Our VFS API has significant differences from upstream;
follow the [usage example](_examples/library/vfs/main.go).

Our `PRAGMA litestream_time` accepts:
- Go [duration strings](https://pkg.go.dev/time#ParseDuration)
- SQLite [time values](https://sqlite.org/lang_datefunc.html#time_values)
- SQLite [time modifiers 1 through 13](https://sqlite.org/lang_datefunc.html#modifiers)
