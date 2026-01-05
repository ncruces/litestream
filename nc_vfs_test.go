package litestream_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/litestream"
	"github.com/ncruces/litestream/file"
)

func TestCreateVFSReadReplica(t *testing.T) {
	dir := t.TempDir()
	dbpath := filepath.Join(dir, "test.db")
	backup := filepath.Join(dir, "backup", "test.db")

	db, err := driver.Open(dbpath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	client := file.NewReplicaClient(backup)
	setupVFSReplication(t, dbpath, client)

	litestream.CreateVFSReadReplica("test.db", client, litestream.VFSOptions{})
	replica, err := driver.Open("file:test.db?vfs=litestream")
	if err != nil {
		t.Fatal(err)
	}
	defer replica.Close()

	_, err = db.ExecContext(t.Context(), `CREATE TABLE users (id INT, name VARCHAR(10))`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.ExecContext(t.Context(),
		`INSERT INTO users (id, name) VALUES (0, 'go'), (1, 'zig'), (2, 'whatever')`)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(litestream.DefaultPollInterval + litestream.DefaultMonitorInterval)

	rows, err := replica.QueryContext(t.Context(), `SELECT id, name FROM users`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	row := 0
	ids := []int{0, 1, 2}
	names := []string{"go", "zig", "whatever"}
	for ; rows.Next(); row++ {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			t.Fatal(err)
		}

		if id != ids[row] {
			t.Errorf("got %d, want %d", id, ids[row])
		}
		if name != names[row] {
			t.Errorf("got %q, want %q", name, names[row])
		}
	}
	if row != 3 {
		t.Errorf("got %d, want %d", row, len(ids))
	}

	var lag int
	err = replica.QueryRowContext(t.Context(), `PRAGMA litestream_lag`).Scan(&lag)
	if err != nil {
		t.Fatal(err)
	}
	if lag < 0 || lag > 2 {
		t.Errorf("got %d", lag)
	}

	var txid string
	err = replica.QueryRowContext(t.Context(), `PRAGMA litestream_txid`).Scan(&txid)
	if err != nil {
		t.Fatal(err)
	}
	if txid != "0000000000000001" {
		t.Errorf("got %q", txid)
	}
}

func setupVFSReplication(tb testing.TB, path string, client litestream.ReplicaClient) {
	lsdb := litestream.NewDB(path)
	lsdb.Replica = litestream.NewReplicaWithClient(lsdb, client)

	err := lsdb.Open()
	if err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() { lsdb.Close(tb.Context()) })
}
