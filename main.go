package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/lkarlslund/jetkvm-native/internal/app"
)

func main() {
	var cfg app.Config
	flag.StringVar(&cfg.BaseURL, "base-url", "http://127.0.0.1:8080", "JetKVM device or emulator base URL")
	flag.StringVar(&cfg.Password, "password", "", "Password for local auth mode")
	flag.DurationVar(&cfg.RPCTimeout, "rpc-timeout", 5*time.Second, "Timeout for JSON-RPC requests")
	flag.Parse()

	clientApp, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := clientApp.Connect(ctx); err != nil && !errors.Is(err, context.Canceled) {
			clientApp.SetStatus(fmt.Sprintf("connect failed: %v", err))
		}
	}()

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("jetkvm-native")
	if err := ebiten.RunGame(clientApp); err != nil {
		log.Fatal(err)
	}
}
