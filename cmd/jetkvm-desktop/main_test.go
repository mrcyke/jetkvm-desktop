package main

import (
	"strings"
	"testing"
)

func TestReadPasswordTrimsTrailingLineEndings(t *testing.T) {
	password, err := readPassword(strings.NewReader("secret\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	if password != "secret" {
		t.Fatalf("password = %q, want %q", password, "secret")
	}
}

func TestResolvePasswordPrefersStdin(t *testing.T) {
	password, err := resolvePassword(true, "", strings.NewReader("stdin-secret\n"), func(string) string {
		return "env-secret"
	})
	if err != nil {
		t.Fatal(err)
	}
	if password != "stdin-secret" {
		t.Fatalf("password = %q, want %q", password, "stdin-secret")
	}
}

func TestResolvePasswordUsesNamedEnv(t *testing.T) {
	password, err := resolvePassword(false, "JETKVM_ALT_PASSWORD", strings.NewReader("ignored"), func(name string) string {
		if name != "JETKVM_ALT_PASSWORD" {
			t.Fatalf("env name = %q, want %q", name, "JETKVM_ALT_PASSWORD")
		}
		return "named-secret"
	})
	if err != nil {
		t.Fatal(err)
	}
	if password != "named-secret" {
		t.Fatalf("password = %q, want %q", password, "named-secret")
	}
}

func TestResolvePasswordFallsBackToDefaultEnv(t *testing.T) {
	password, err := resolvePassword(false, "", strings.NewReader("ignored"), func(name string) string {
		if name != defaultPasswordEnv {
			t.Fatalf("env name = %q, want %q", name, defaultPasswordEnv)
		}
		return "default-secret"
	})
	if err != nil {
		t.Fatal(err)
	}
	if password != "default-secret" {
		t.Fatalf("password = %q, want %q", password, "default-secret")
	}
}

func TestResolvePasswordRejectsConflictingSources(t *testing.T) {
	_, err := resolvePassword(true, "JETKVM_ALT_PASSWORD", strings.NewReader("stdin-secret"), func(string) string {
		return "env-secret"
	})
	if err == nil {
		t.Fatal("expected conflict error")
	}
}
