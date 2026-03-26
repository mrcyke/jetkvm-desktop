package emulator

import (
	"context"
	"testing"
	"time"

	"github.com/lkarlslund/jetkvm-native/internal/client"
)

func TestClientConnectsAndRPCWorks(t *testing.T) {
	srv, err := NewServer(Config{
		ListenAddr: "127.0.0.1:0",
		AuthMode:   AuthModePassword,
		Password:   "secret",
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := srv.ListenAndServe(ctx); err != nil {
			t.Errorf("server: %v", err)
		}
	}()

	deadline := time.Now().Add(2 * time.Second)
	for srv.BaseURL() == "" && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if srv.BaseURL() == "" {
		t.Fatal("server did not start")
	}

	c, err := client.New(client.Config{
		BaseURL:    srv.BaseURL(),
		Password:   "secret",
		RPCTimeout: 2 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if err := c.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	waitCtx, waitCancel := context.WithTimeout(ctx, 5*time.Second)
	defer waitCancel()
	if err := c.WaitForHID(waitCtx); err != nil {
		t.Fatal(err)
	}

	var pong string
	if err := c.Call(waitCtx, "ping", nil, &pong); err != nil {
		t.Fatal(err)
	}
	if pong != "pong" {
		t.Fatalf("expected pong, got %q", pong)
	}

	if err := c.SendKeypress(4, true); err != nil {
		t.Fatal(err)
	}

	deadline = time.Now().Add(2 * time.Second)
	for len(srv.Inputs()) == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if len(srv.Inputs()) == 0 {
		t.Fatal("expected HID input to be recorded")
	}
}
