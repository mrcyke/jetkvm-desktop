package input

import (
	"testing"
)

func TestKeyboardUpdatePressAndRelease(t *testing.T) {
	k := NewKeyboard()

	events := k.Update([]Key{KeyA, KeyShiftLeft})
	if len(events) != 2 {
		t.Fatalf("expected 2 press events, got %d", len(events))
	}
	if events[0] != (KeyEvent{HID: 4, Press: true}) {
		t.Fatalf("unexpected first event: %+v", events[0])
	}
	if events[1] != (KeyEvent{HID: 225, Press: true}) {
		t.Fatalf("unexpected second event: %+v", events[1])
	}

	events = k.Update([]Key{KeyShiftLeft})
	if len(events) != 1 {
		t.Fatalf("expected 1 release event, got %d", len(events))
	}
	if events[0] != (KeyEvent{HID: 4, Press: false}) {
		t.Fatalf("unexpected release event: %+v", events[0])
	}
}

func TestKeyboardReleaseAll(t *testing.T) {
	k := NewKeyboard()
	_ = k.Update([]Key{KeyB, KeyControlRight})

	events := k.ReleaseAll()
	if len(events) != 2 {
		t.Fatalf("expected 2 release events, got %d", len(events))
	}
	if events[0].Press || events[1].Press {
		t.Fatalf("expected release events only, got %+v", events)
	}

	events = k.ReleaseAll()
	if len(events) != 0 {
		t.Fatalf("expected no events after repeated release all, got %+v", events)
	}
}

func TestKeyboardIgnoresUnknownKeys(t *testing.T) {
	k := NewKeyboard()

	events := k.Update([]Key{KeyUnknown})
	if len(events) != 0 {
		t.Fatalf("expected unknown key to be ignored, got %+v", events)
	}
}
