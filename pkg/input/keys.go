package input

type Key int

const (
	KeyUnknown Key = iota
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	Key0
	KeyEnter
	KeyEscape
	KeyBackspace
	KeyTab
	KeySpace
	KeyMinus
	KeyEqual
	KeyLeftBracket
	KeyRightBracket
	KeyBackslash
	KeySemicolon
	KeyApostrophe
	KeyGraveAccent
	KeyComma
	KeyPeriod
	KeySlash
	KeyCapsLock
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyPrintScreen
	KeyScrollLock
	KeyPause
	KeyInsert
	KeyHome
	KeyPageUp
	KeyDelete
	KeyEnd
	KeyPageDown
	KeyRight
	KeyLeft
	KeyDown
	KeyUp
	KeyNumLock
	KeyNumpadDivide
	KeyNumpadMultiply
	KeyNumpadSubtract
	KeyNumpadAdd
	KeyNumpadEnter
	KeyNumpad1
	KeyNumpad2
	KeyNumpad3
	KeyNumpad4
	KeyNumpad5
	KeyNumpad6
	KeyNumpad7
	KeyNumpad8
	KeyNumpad9
	KeyNumpad0
	KeyNumpadDecimal
	KeyControlLeft
	KeyShiftLeft
	KeyAltLeft
	KeyMetaLeft
	KeyControlRight
	KeyShiftRight
	KeyAltRight
	KeyMetaRight
)

func KeyToHID(key Key) (byte, bool) {
	switch key {
	case KeyA:
		return 4, true
	case KeyB:
		return 5, true
	case KeyC:
		return 6, true
	case KeyD:
		return 7, true
	case KeyE:
		return 8, true
	case KeyF:
		return 9, true
	case KeyG:
		return 10, true
	case KeyH:
		return 11, true
	case KeyI:
		return 12, true
	case KeyJ:
		return 13, true
	case KeyK:
		return 14, true
	case KeyL:
		return 15, true
	case KeyM:
		return 16, true
	case KeyN:
		return 17, true
	case KeyO:
		return 18, true
	case KeyP:
		return 19, true
	case KeyQ:
		return 20, true
	case KeyR:
		return 21, true
	case KeyS:
		return 22, true
	case KeyT:
		return 23, true
	case KeyU:
		return 24, true
	case KeyV:
		return 25, true
	case KeyW:
		return 26, true
	case KeyX:
		return 27, true
	case KeyY:
		return 28, true
	case KeyZ:
		return 29, true
	case Key1:
		return 30, true
	case Key2:
		return 31, true
	case Key3:
		return 32, true
	case Key4:
		return 33, true
	case Key5:
		return 34, true
	case Key6:
		return 35, true
	case Key7:
		return 36, true
	case Key8:
		return 37, true
	case Key9:
		return 38, true
	case Key0:
		return 39, true
	case KeyEnter:
		return 40, true
	case KeyEscape:
		return 41, true
	case KeyBackspace:
		return 42, true
	case KeyTab:
		return 43, true
	case KeySpace:
		return 44, true
	case KeyMinus:
		return 45, true
	case KeyEqual:
		return 46, true
	case KeyLeftBracket:
		return 47, true
	case KeyRightBracket:
		return 48, true
	case KeyBackslash:
		return 49, true
	case KeySemicolon:
		return 51, true
	case KeyApostrophe:
		return 52, true
	case KeyGraveAccent:
		return 53, true
	case KeyComma:
		return 54, true
	case KeyPeriod:
		return 55, true
	case KeySlash:
		return 56, true
	case KeyCapsLock:
		return 57, true
	case KeyF1:
		return 58, true
	case KeyF2:
		return 59, true
	case KeyF3:
		return 60, true
	case KeyF4:
		return 61, true
	case KeyF5:
		return 62, true
	case KeyF6:
		return 63, true
	case KeyF7:
		return 64, true
	case KeyF8:
		return 65, true
	case KeyF9:
		return 66, true
	case KeyF10:
		return 67, true
	case KeyF11:
		return 68, true
	case KeyF12:
		return 69, true
	case KeyPrintScreen:
		return 70, true
	case KeyScrollLock:
		return 71, true
	case KeyPause:
		return 72, true
	case KeyInsert:
		return 73, true
	case KeyHome:
		return 74, true
	case KeyPageUp:
		return 75, true
	case KeyDelete:
		return 76, true
	case KeyEnd:
		return 77, true
	case KeyPageDown:
		return 78, true
	case KeyRight:
		return 79, true
	case KeyLeft:
		return 80, true
	case KeyDown:
		return 81, true
	case KeyUp:
		return 82, true
	case KeyNumLock:
		return 83, true
	case KeyNumpadDivide:
		return 84, true
	case KeyNumpadMultiply:
		return 85, true
	case KeyNumpadSubtract:
		return 86, true
	case KeyNumpadAdd:
		return 87, true
	case KeyNumpadEnter:
		return 88, true
	case KeyNumpad1:
		return 89, true
	case KeyNumpad2:
		return 90, true
	case KeyNumpad3:
		return 91, true
	case KeyNumpad4:
		return 92, true
	case KeyNumpad5:
		return 93, true
	case KeyNumpad6:
		return 94, true
	case KeyNumpad7:
		return 95, true
	case KeyNumpad8:
		return 96, true
	case KeyNumpad9:
		return 97, true
	case KeyNumpad0:
		return 98, true
	case KeyNumpadDecimal:
		return 99, true
	case KeyControlLeft:
		return 224, true
	case KeyShiftLeft:
		return 225, true
	case KeyAltLeft:
		return 226, true
	case KeyMetaLeft:
		return 227, true
	case KeyControlRight:
		return 228, true
	case KeyShiftRight:
		return 229, true
	case KeyAltRight:
		return 230, true
	case KeyMetaRight:
		return 231, true
	default:
		return 0, false
	}
}
