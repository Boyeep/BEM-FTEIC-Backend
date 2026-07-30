package middleware

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func TestRedisRateLimitStoreIncrementsAndExpires(t *testing.T) {
	server := miniredis.RunT(t)
	store, err := NewRedisRateLimitStore(
		fmt.Sprintf("redis://%s/0", server.Addr()),
		"test",
	)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	count, reset, err := store.Increment("client", now, time.Minute)
	if err != nil || count != 1 || !reset.After(now) {
		t.Fatalf("unexpected first increment: count=%d reset=%v err=%v", count, reset, err)
	}
	count, _, err = store.Increment("client", now, time.Minute)
	if err != nil || count != 2 {
		t.Fatalf("unexpected second increment: count=%d err=%v", count, err)
	}
	server.FastForward(time.Minute)
	count, _, err = store.Increment("client", now.Add(time.Minute), time.Minute)
	if err != nil || count != 1 {
		t.Fatalf("expected reset count, got %d: %v", count, err)
	}
}
