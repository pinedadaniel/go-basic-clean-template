package config

import (
	"testing"
)

func setRequiredEnvironment(t *testing.T) {
	t.Helper()

	t.Setenv("APP_NAME", "test-app")
	t.Setenv("APP_VERSION", "1.2.3")
	t.Setenv("APP_SCOPE", "LoCaL")
	t.Setenv("PROFILE_DIR", t.TempDir())
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_FORMAT", "json")
}

func TestNewLoadsEnvironmentConfig(t *testing.T) {
	setRequiredEnvironment(t)

	got, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if got.App.Name != "test-app" {
		t.Errorf("App.Name = %q, want %q", got.App.Name, "test-app")
	}
	if got.App.Scope != "LoCaL" {
		t.Errorf("App.Scope = %q, want %q", got.App.Scope, "LoCaL")
	}
	if got.ProfileDir == "" {
		t.Error("ProfileDir is empty")
	}
}

func TestNewRejectsInvalidScope(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("APP_SCOPE", "staging")

	if _, err := New(); err == nil {
		t.Fatal("New() error = nil, want invalid scope error")
	}
}

func TestNewRequiresProfileDir(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("PROFILE_DIR", "")

	if _, err := New(); err == nil {
		t.Fatal("New() error = nil, want missing PROFILE_DIR error")
	}
}
