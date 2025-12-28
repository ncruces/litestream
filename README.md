# Litestream for the `ncruces/go-sqlite3` driver

See: [github.com/benbjohnson/litestream](https://github.com/benbjohnson/litestream)

This is a fork focusing on library use, which uses the
[github.com/ncruces/go-sqlite3](https://github.com/ncruces/go-sqlite3)
driver.

It implements the `"litestream"` SQLite VFS that offers
Litestream [lightweight read-replicas](https://fly.io/blog/litestream-revamped/#lightweight-read-replicas).

Our `PRAGMA litestream_time` accepts:
- Go [duration strings](https://pkg.go.dev/time#ParseDuration)
- SQLite [time values](https://sqlite.org/lang_datefunc.html#time_values)
- SQLite [time modifiers 1 through 13](https://sqlite.org/lang_datefunc.html#modifiers)

Checkout the [examples](_examples/library/) for how to use.
