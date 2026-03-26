package app

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"sort"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/lkarlslund/jetkvm-native/pkg/session"
)

type Config struct {
	BaseURL    string
	Password   string
	RPCTimeout time.Duration
}

type App struct {
	cfg  Config
	ctrl *session.Controller

	mu         sync.RWMutex
	lastImg    *ebiten.Image
	keys       map[ebiten.Key]bool
	lastX      int
	lastY      int
	relative   bool
	renderRect rect
	focused    bool
}

func New(cfg Config) (*App, error) {
	ctrl := session.New(session.Config{
		BaseURL:    cfg.BaseURL,
		Password:   cfg.Password,
		RPCTimeout: cfg.RPCTimeout,
		Reconnect:  true,
	})
	return &App{
		cfg:     cfg,
		ctrl:    ctrl,
		keys:    map[ebiten.Key]bool{},
		focused: true,
	}, nil
}

func (a *App) Start(ctx context.Context) {
	a.ctrl.Start(ctx)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			frame := a.ctrl.LatestFrame()
			if frame == nil {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			img := ebiten.NewImageFromImage(frame)
			a.mu.Lock()
			a.lastImg = img
			a.mu.Unlock()
			time.Sleep(16 * time.Millisecond)
		}
	}()
}

func (a *App) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	nowFocused := ebiten.IsFocused()
	if a.focused && !nowFocused {
		a.releaseAllKeys()
		if a.relative {
			a.relative = false
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		}
	}
	a.focused = nowFocused

	if inpututil.IsKeyJustPressed(ebiten.KeyF8) {
		a.relative = !a.relative
		if a.relative {
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		} else {
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		}
		a.lastX, a.lastY = ebiten.CursorPosition()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		_ = a.ctrl.Reboot()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		a.adjustStreamQuality(+0.05)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		a.adjustStreamQuality(-0.05)
	}

	a.syncKeyboard()
	a.syncMouse()
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	snap := a.ctrl.Snapshot()
	screen.Fill(color.RGBA{R: 10, G: 15, B: 24, A: 255})
	vector.DrawFilledRect(screen, 0, 0, float32(screen.Bounds().Dx()), 44, color.RGBA{R: 16, G: 28, B: 44, A: 255}, false)
	vector.DrawFilledRect(screen, 0, float32(screen.Bounds().Dy()-32), float32(screen.Bounds().Dx()), 32, color.RGBA{R: 16, G: 28, B: 44, A: 255}, false)

	videoArea := image.Rect(16, 56, screen.Bounds().Dx()-16, screen.Bounds().Dy()-44)
	a.mu.RLock()
	img := a.lastImg
	a.mu.RUnlock()
	if img != nil {
		w, h := img.Bounds().Dx(), img.Bounds().Dy()
		op := &ebiten.DrawImageOptions{}
		scale := min(float64(videoArea.Dx())/float64(w), float64(videoArea.Dy())/float64(h))
		drawW := float64(w) * scale
		drawH := float64(h) * scale
		x := float64(videoArea.Min.X) + (float64(videoArea.Dx())-drawW)/2
		y := float64(videoArea.Min.Y) + (float64(videoArea.Dy())-drawH)/2
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x, y)
		screen.DrawImage(img, op)
		a.renderRect = rect{
			x: x,
			y: y,
			w: drawW,
			h: drawH,
		}
	} else {
		a.renderRect = rect{}
	}
	mode := "absolute"
	if a.relative {
		mode = "relative"
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("jetkvm-client  %s", snap.BaseURL), 16, 14)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("phase: %s  hid: %t  rtc: %s  quality: %.2f  mouse: %s", snap.Phase, snap.HIDReady, snap.RTCState, snap.Quality, mode), 16, screen.Bounds().Dy()-24)
	if snap.DeviceID != "" || snap.Hostname != "" {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s  %s", snap.DeviceID, snap.Hostname), screen.Bounds().Dx()-320, 14)
	}
	a.drawOverlay(screen, snap, img != nil)
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (a *App) syncKeyboard() {
	if a.ctrl.Snapshot().Phase != session.PhaseConnected {
		return
	}
	keys := inpututil.AppendPressedKeys(nil)
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	current := make(map[ebiten.Key]bool, len(keys))
	for _, key := range keys {
		current[key] = true
		if !a.keys[key] {
			if hid, ok := keyToHID(key); ok {
				_ = a.ctrl.SendKeypress(hid, true)
			}
		}
	}
	for key := range a.keys {
		if current[key] {
			continue
		}
		if hid, ok := keyToHID(key); ok {
			_ = a.ctrl.SendKeypress(hid, false)
		}
	}
	a.keys = current
}

func (a *App) syncMouse() {
	if a.ctrl.Snapshot().Phase != session.PhaseConnected {
		return
	}
	x, y := ebiten.CursorPosition()
	buttons := byte(0)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		buttons |= 1
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		buttons |= 2
	}
	if a.relative {
		dx := int8(clamp(float64(x-a.lastX), -127, 127))
		dy := int8(clamp(float64(y-a.lastY), -127, 127))
		if dx != 0 || dy != 0 || buttons != 0 {
			_ = a.ctrl.SendRelMouse(dx, dy, buttons)
		}
	} else {
		if !a.renderRect.valid() {
			return
		}
		nx, ny := a.renderRect.toHID(x, y)
		_ = a.ctrl.SendAbsPointer(nx, ny, buttons)
	}
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		_ = a.ctrl.SendWheel(int8(clamp(-wheelY, -127, 127)))
	}
	a.lastX = x
	a.lastY = y
}

func keyToHID(key ebiten.Key) (byte, bool) {
	switch key {
	case ebiten.KeyA:
		return 4, true
	case ebiten.KeyB:
		return 5, true
	case ebiten.KeyC:
		return 6, true
	case ebiten.KeyD:
		return 7, true
	case ebiten.KeyE:
		return 8, true
	case ebiten.KeyF:
		return 9, true
	case ebiten.KeyG:
		return 10, true
	case ebiten.KeyH:
		return 11, true
	case ebiten.KeyI:
		return 12, true
	case ebiten.KeyJ:
		return 13, true
	case ebiten.KeyK:
		return 14, true
	case ebiten.KeyL:
		return 15, true
	case ebiten.KeyM:
		return 16, true
	case ebiten.KeyN:
		return 17, true
	case ebiten.KeyO:
		return 18, true
	case ebiten.KeyP:
		return 19, true
	case ebiten.KeyQ:
		return 20, true
	case ebiten.KeyR:
		return 21, true
	case ebiten.KeyS:
		return 22, true
	case ebiten.KeyT:
		return 23, true
	case ebiten.KeyU:
		return 24, true
	case ebiten.KeyV:
		return 25, true
	case ebiten.KeyW:
		return 26, true
	case ebiten.KeyX:
		return 27, true
	case ebiten.KeyY:
		return 28, true
	case ebiten.KeyZ:
		return 29, true
	case ebiten.Key1:
		return 30, true
	case ebiten.Key2:
		return 31, true
	case ebiten.Key3:
		return 32, true
	case ebiten.Key4:
		return 33, true
	case ebiten.Key5:
		return 34, true
	case ebiten.Key6:
		return 35, true
	case ebiten.Key7:
		return 36, true
	case ebiten.Key8:
		return 37, true
	case ebiten.Key9:
		return 38, true
	case ebiten.Key0:
		return 39, true
	case ebiten.KeyEnter:
		return 40, true
	case ebiten.KeyEscape:
		return 41, true
	case ebiten.KeyBackspace:
		return 42, true
	case ebiten.KeyTab:
		return 43, true
	case ebiten.KeySpace:
		return 44, true
	case ebiten.KeyMinus:
		return 45, true
	case ebiten.KeyEqual:
		return 46, true
	case ebiten.KeyLeftBracket:
		return 47, true
	case ebiten.KeyRightBracket:
		return 48, true
	case ebiten.KeyBackslash:
		return 49, true
	case ebiten.KeySemicolon:
		return 51, true
	case ebiten.KeyApostrophe:
		return 52, true
	case ebiten.KeyGraveAccent:
		return 53, true
	case ebiten.KeyComma:
		return 54, true
	case ebiten.KeyPeriod:
		return 55, true
	case ebiten.KeySlash:
		return 56, true
	case ebiten.KeyCapsLock:
		return 57, true
	case ebiten.KeyF1:
		return 58, true
	case ebiten.KeyF2:
		return 59, true
	case ebiten.KeyF3:
		return 60, true
	case ebiten.KeyF4:
		return 61, true
	case ebiten.KeyF5:
		return 62, true
	case ebiten.KeyF6:
		return 63, true
	case ebiten.KeyF7:
		return 64, true
	case ebiten.KeyF8:
		return 65, true
	case ebiten.KeyF9:
		return 66, true
	case ebiten.KeyF10:
		return 67, true
	case ebiten.KeyF11:
		return 68, true
	case ebiten.KeyF12:
		return 69, true
	case ebiten.KeyPrintScreen:
		return 70, true
	case ebiten.KeyScrollLock:
		return 71, true
	case ebiten.KeyPause:
		return 72, true
	case ebiten.KeyInsert:
		return 73, true
	case ebiten.KeyHome:
		return 74, true
	case ebiten.KeyPageUp:
		return 75, true
	case ebiten.KeyDelete:
		return 76, true
	case ebiten.KeyEnd:
		return 77, true
	case ebiten.KeyPageDown:
		return 78, true
	case ebiten.KeyRight:
		return 79, true
	case ebiten.KeyLeft:
		return 80, true
	case ebiten.KeyDown:
		return 81, true
	case ebiten.KeyUp:
		return 82, true
	case ebiten.KeyNumLock:
		return 83, true
	case ebiten.KeyNumpadDivide:
		return 84, true
	case ebiten.KeyNumpadMultiply:
		return 85, true
	case ebiten.KeyNumpadSubtract:
		return 86, true
	case ebiten.KeyNumpadAdd:
		return 87, true
	case ebiten.KeyNumpadEnter:
		return 88, true
	case ebiten.KeyNumpad1:
		return 89, true
	case ebiten.KeyNumpad2:
		return 90, true
	case ebiten.KeyNumpad3:
		return 91, true
	case ebiten.KeyNumpad4:
		return 92, true
	case ebiten.KeyNumpad5:
		return 93, true
	case ebiten.KeyNumpad6:
		return 94, true
	case ebiten.KeyNumpad7:
		return 95, true
	case ebiten.KeyNumpad8:
		return 96, true
	case ebiten.KeyNumpad9:
		return 97, true
	case ebiten.KeyNumpad0:
		return 98, true
	case ebiten.KeyNumpadDecimal:
		return 99, true
	case ebiten.KeyControlLeft:
		return 224, true
	case ebiten.KeyShiftLeft:
		return 225, true
	case ebiten.KeyAltLeft:
		return 226, true
	case ebiten.KeyMetaLeft:
		return 227, true
	case ebiten.KeyControlRight:
		return 228, true
	case ebiten.KeyShiftRight:
		return 229, true
	case ebiten.KeyAltRight:
		return 230, true
	case ebiten.KeyMetaRight:
		return 231, true
	default:
		return 0, false
	}
}

func clamp(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (a *App) adjustStreamQuality(delta float64) {
	snap := a.ctrl.Snapshot()
	next := clamp(snap.Quality+delta, 0.1, 1.0)

	go func(value float64) {
		_ = a.ctrl.SetQuality(value)
	}(next)
}

func (a *App) releaseAllKeys() {
	for key := range a.keys {
		if hid, ok := keyToHID(key); ok {
			_ = a.ctrl.SendKeypress(hid, false)
		}
	}
	a.keys = map[ebiten.Key]bool{}
}

func (a *App) drawOverlay(screen *ebiten.Image, snap session.Snapshot, hasVideo bool) {
	message := ""
	switch snap.Phase {
	case session.PhaseConnecting:
		message = "Connecting to device"
	case session.PhaseReconnecting:
		message = "Reconnecting to device"
	case session.PhaseAuthFailed:
		message = "Authentication failed"
	case session.PhaseOtherSession:
		message = "Another session took over"
	case session.PhaseRebooting:
		message = "Device rebooting"
	case session.PhaseDisconnected:
		message = "Connection lost"
	case session.PhaseFatal:
		message = "Fatal error"
	case session.PhaseConnected:
		if !hasVideo || !snap.VideoReady {
			message = "Loading video stream"
		}
	}
	if message == "" {
		return
	}
	vector.DrawFilledRect(screen, 24, 72, float32(screen.Bounds().Dx()-48), 72, color.RGBA{R: 8, G: 12, B: 18, A: 224}, false)
	ebitenutil.DebugPrintAt(screen, message, 40, 96)
	if snap.LastError != "" && snap.Phase != session.PhaseConnected {
		ebitenutil.DebugPrintAt(screen, snap.LastError, 40, 116)
	}
}

type rect struct {
	x float64
	y float64
	w float64
	h float64
}

func (r rect) valid() bool {
	return r.w > 0 && r.h > 0
}

func (r rect) toHID(cursorX, cursorY int) (uint16, uint16) {
	if !r.valid() {
		return 0, 0
	}
	relX := clamp((float64(cursorX)-r.x)/r.w, 0, 1)
	relY := clamp((float64(cursorY)-r.y)/r.h, 0, 1)
	return uint16(relX * 32767.0), uint16(relY * 32767.0)
}
