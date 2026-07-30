package storageredis

import (
	"runtime"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// The library default (5 x GOMAXPROCS) collapsed in production: an on-demand
// TLS edge parks thousands of concurrent handshakes behind it. The default
// here must be sized for that workload without requiring configuration.
func TestPoolDefaultsAreSizedForOnDemandTLS(t *testing.T) {
	rs := &RedisStorage{}
	opts := redis.UniversalOptions{}
	rs.applyPoolConfig(&opts)

	want := 50 * runtime.GOMAXPROCS(0)
	if opts.PoolSize != want {
		t.Errorf("PoolSize = %d, want %d (50x GOMAXPROCS)", opts.PoolSize, want)
	}
	if opts.MinIdleConns != 5 {
		t.Errorf("MinIdleConns = %d, want 5 pre-warmed", opts.MinIdleConns)
	}
	if opts.PoolTimeout != 0 {
		t.Errorf("PoolTimeout = %v, want 0 (library default) when unconfigured", opts.PoolTimeout)
	}
}

func TestPoolConfigIsHonoredWhenSet(t *testing.T) {
	rs := &RedisStorage{PoolSize: 300, MinIdleConns: 20, PoolTimeout: "2"}
	opts := redis.UniversalOptions{}
	rs.applyPoolConfig(&opts)

	if opts.PoolSize != 300 || opts.MinIdleConns != 20 {
		t.Errorf("explicit pool settings not honored: %+v", opts)
	}
	if opts.PoolTimeout != 2*time.Second {
		t.Errorf("PoolTimeout = %v, want 2s", opts.PoolTimeout)
	}
}

func TestPoolTimeoutIgnoresGarbage(t *testing.T) {
	rs := &RedisStorage{PoolTimeout: "nonsense"}
	opts := redis.UniversalOptions{}
	rs.applyPoolConfig(&opts)

	if opts.PoolTimeout != 0 {
		t.Errorf("PoolTimeout = %v, want untouched on unparseable input", opts.PoolTimeout)
	}
}
