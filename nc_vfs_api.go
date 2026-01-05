package litestream

import (
	"log/slog"
	"sync"
	"time"

	"github.com/ncruces/go-sqlite3/vfs"
)

const (
	// The default poll interval.
	DefaultPollInterval = 1 * time.Second

	// The default cache size: 10 MiB.
	DefaultCacheSize = 10 * 1024 * 1024
)

func init() {
	vfs.Register("litestream", liteVFS{})
}

var (
	liteMtx sync.RWMutex
	// +checklocks:liteMtx
	liteDBs = map[string]*liteDB{}
)

// VFSOptions represents options for [CreateVFSReadReplica].
type VFSOptions struct {
	// Where to log error messages. May be nil.
	Logger *slog.Logger

	// Replica poll interval.
	// Should be less than the compaction interval
	// used by the [ReplicaClient] at MinLevel+1.
	PollInterval time.Duration

	// CacheSize is the maximum size for the page cache in bytes.
	// Zero means [DefaultCacheSize], negative disables caching.
	CacheSize int
}

// CreateVFSReadReplica creates a read-replica VFS from a client.
func CreateVFSReadReplica(name string, client ReplicaClient, options VFSOptions) {
	if options.Logger != nil {
		options.Logger = options.Logger.With("name", name)
	} else {
		options.Logger = slog.New(slog.DiscardHandler)
	}
	if options.PollInterval <= 0 {
		options.PollInterval = DefaultPollInterval
	}
	if options.CacheSize == 0 {
		options.CacheSize = DefaultCacheSize
	}

	liteMtx.Lock()
	liteDBs[name] = &liteDB{
		client: client,
		opts:   options,
		cache:  pageCache{size: options.CacheSize},
	}
	liteMtx.Unlock()
}

// RemoveVFSReadReplica removes a Litestream read-replica VFS by name.
func RemoveVFSReadReplica(name string) {
	liteMtx.Lock()
	delete(liteDBs, name)
	liteMtx.Unlock()
}
