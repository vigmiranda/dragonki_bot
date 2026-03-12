package integrationtest

import (
	"context"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisInstance struct {
	URL    string
	Client *redis.Client
}

func NewRedisInstance(t *testing.T) *RedisInstance {
	t.Helper()

	rawURL := strings.TrimSpace(os.Getenv("TEST_REDIS_URL"))
	if rawURL == "" {
		t.Skip("set TEST_REDIS_URL to run redis integration tests")
	}
	MustUseSafeRedisURL(t, rawURL)

	opts, err := redis.ParseURL(rawURL)
	if err != nil {
		t.Fatalf("parse TEST_REDIS_URL: %v", err)
	}

	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		t.Fatalf("ping redis: %v", err)
	}

	if err := client.FlushDB(ctx).Err(); err != nil {
		_ = client.Close()
		t.Fatalf("flush redis db before test: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		_ = client.FlushDB(cleanupCtx).Err()
		_ = client.Close()
	})

	return &RedisInstance{
		URL:    rawURL,
		Client: client,
	}
}

func MustUseSafeRedisURL(t *testing.T, rawURL string) {
	t.Helper()

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("parse TEST_REDIS_URL: %v", err)
	}

	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	if host == "" {
		t.Fatal("TEST_REDIS_URL must include a hostname")
	}

	allowedHosts := map[string]struct{}{
		"localhost":      {},
		"127.0.0.1":      {},
		"::1":            {},
		"redis":          {},
		"dragonki_redis": {},
	}

	if _, ok := allowedHosts[host]; ok {
		return
	}
	if strings.HasSuffix(host, ".local") || strings.HasSuffix(host, ".internal") {
		return
	}

	t.Fatalf("refusing to run integration tests against non-local TEST_REDIS_URL host %q", host)
}
