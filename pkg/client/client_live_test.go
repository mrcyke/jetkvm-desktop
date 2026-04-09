package client

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"
)

func TestLiveDeviceConnectsAndStreams(t *testing.T) {
	baseURL := os.Getenv("JETKVM_BASE_URL")
	if baseURL == "" {
		t.Skip("JETKVM_BASE_URL not set")
	}

	password := os.Getenv("JETKVM_PASSWORD")

	c, err := New(Config{
		BaseURL:    baseURL,
		Password:   password,
		RPCTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	var (
		mu     sync.Mutex
		events []LifecycleEvent
	)
	logCtx, logCancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-logCtx.Done():
				return
			case evt := <-c.Lifecycle():
				mu.Lock()
				events = append(events, evt)
				mu.Unlock()
			}
		}
	}()
	defer func() {
		logCancel()
		select {
		case <-done:
		case <-time.After(500 * time.Millisecond):
		}
		mu.Lock()
		defer mu.Unlock()
		for _, evt := range events {
			t.Logf("lifecycle: type=%s state=%s err=%s", evt.Type, evt.Connection.String(), evt.Err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := c.Connect(ctx); err != nil {
		t.Fatal(err)
	}
	if err := c.WaitForHID(ctx); err != nil {
		t.Fatal(err)
	}

	var deviceID string
	if err := c.Call(ctx, "getDeviceID", nil, &deviceID); err != nil {
		t.Fatal(err)
	}
	if deviceID == "" {
		t.Fatal("expected non-empty device ID")
	}

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		stream := c.VideoStream()
		if stream != nil && stream.Latest() != nil {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}

	if stream := c.VideoStream(); stream != nil && stream.Err() != nil {
		t.Logf("video error: %v", stream.Err())
	}

	t.Fatal("expected at least one decoded video frame from live device")
}
