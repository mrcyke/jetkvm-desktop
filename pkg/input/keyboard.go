package input

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type KeyEvent struct {
	HID   byte
	Press bool
}

type Keyboard struct {
	pressed map[ebiten.Key]bool
}

func NewKeyboard() *Keyboard {
	return &Keyboard{
		pressed: map[ebiten.Key]bool{},
	}
}

func (k *Keyboard) Update(keys []ebiten.Key) []KeyEvent {
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	current := make(map[ebiten.Key]bool, len(keys))
	events := make([]KeyEvent, 0, len(keys)+len(k.pressed))

	for _, key := range keys {
		current[key] = true
		if k.pressed[key] {
			continue
		}
		if hid, ok := KeyToHID(key); ok {
			events = append(events, KeyEvent{HID: hid, Press: true})
		}
	}

	for key := range k.pressed {
		if current[key] {
			continue
		}
		if hid, ok := KeyToHID(key); ok {
			events = append(events, KeyEvent{HID: hid, Press: false})
		}
	}

	k.pressed = current
	return events
}

func (k *Keyboard) ReleaseAll() []KeyEvent {
	keys := make([]ebiten.Key, 0, len(k.pressed))
	for key := range k.pressed {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	events := make([]KeyEvent, 0, len(keys))
	for _, key := range keys {
		if hid, ok := KeyToHID(key); ok {
			events = append(events, KeyEvent{HID: hid, Press: false})
		}
	}
	k.pressed = map[ebiten.Key]bool{}
	return events
}

func KeyToHID(key ebiten.Key) (byte, bool) {
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
