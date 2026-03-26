package app

import (
	"context"
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/lkarlslund/jetkvm-native/internal/client"
)

type Config struct {
	BaseURL    string
	Password   string
	RPCTimeout time.Duration
}

type App struct {
	cfg    Config
	client *client.Client

	mu     sync.RWMutex
	status string
}

func New(cfg Config) (*App, error) {
	c, err := client.New(client.Config{
		BaseURL:    cfg.BaseURL,
		Password:   cfg.Password,
		RPCTimeout: cfg.RPCTimeout,
	})
	if err != nil {
		return nil, err
	}
	return &App{
		cfg:    cfg,
		client: c,
		status: "starting",
	}, nil
}

func (a *App) Connect(ctx context.Context) error {
	a.SetStatus("connecting")
	if err := a.client.Connect(ctx); err != nil {
		return err
	}

	waitCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := a.client.WaitForHID(waitCtx); err != nil {
		return err
	}

	var deviceID string
	if err := a.client.Call(ctx, "getDeviceID", nil, &deviceID); err != nil {
		return err
	}
	a.SetStatus(fmt.Sprintf("connected to %s", deviceID))

	go func() {
		for evt := range a.client.Events() {
			a.SetStatus("event: " + evt.Method)
		}
	}()
	return nil
}

func (a *App) SetStatus(status string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.status = status
}

func (a *App) statusLine() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.status
}

func (a *App) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 12, G: 20, B: 32, A: 255})
	ebitenutil.DebugPrint(screen, "jetkvm-native\n"+a.statusLine()+"\nEsc to quit")
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
