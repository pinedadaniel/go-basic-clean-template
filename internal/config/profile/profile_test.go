package profile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pinedadaniel/go-basic-clean-template/pkg/env"
)

func TestNewLoadsProfileForScopeFromConfiguredDirectory(t *testing.T) {
	dir := t.TempDir()
	profileJSON := `{"another_api":{"base_url":"https://example.com","timeout":1000,"circuit_breaker_ratio":0.5}}`

	for _, scope := range []env.Scope{env.Local, env.Beta, env.Prod} {
		t.Run(string(scope), func(t *testing.T) {
			if err := os.WriteFile(filepath.Join(dir, string(scope)+".json"), []byte(profileJSON), 0o600); err != nil {
				t.Fatalf("write profile: %v", err)
			}

			got, err := New(dir, scope)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			if got.AnotherAPI.BaseURL != "https://example.com" {
				t.Errorf("base URL = %q, want %q", got.AnotherAPI.BaseURL, "https://example.com")
			}
		})
	}
}

func TestNewLoadsProfileForMixedCaseScope(t *testing.T) {
	dir := t.TempDir()
	profileJSON := `{"another_api":{"base_url":"https://example.com","timeout":1000,"circuit_breaker_ratio":0.5}}`
	if err := os.WriteFile(filepath.Join(dir, "local.json"), []byte(profileJSON), 0o600); err != nil {
		t.Fatalf("write profile: %v", err)
	}

	got, err := New(dir, env.Scope(strings.ToUpper(string(env.Local))))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got.AnotherAPI.BaseURL != "https://example.com" {
		t.Errorf("base URL = %q, want %q", got.AnotherAPI.BaseURL, "https://example.com")
	}
}

func TestNewLoadsLocalProfileWithStrictDecoder(t *testing.T) {
	got, err := New(".", env.Local)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got.AnotherAPI.Version != "/v1" {
		t.Errorf("version = %q, want %q", got.AnotherAPI.Version, "/v1")
	}
}

func TestProfileValidateRejectsInvalidEndpointConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		wantErr string
	}{
		{
			name:    "malformed URL",
			baseURL: "://example.com",
			wantErr: "another_api.base_url",
		},
		{
			name:    "URL without host",
			baseURL: "https:///path",
			wantErr: "another_api.base_url",
		},
		{
			name:    "unsupported scheme",
			baseURL: "ftp://example.com",
			wantErr: "another_api.base_url",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			profile := Profile{
				AnotherAPI: EndpointConfiguration{
					BaseURL: test.baseURL,
					Timeout: 1000,
				},
			}

			err := profile.validate()
			if err == nil {
				t.Fatal("validate() error = nil")
			}
			if !strings.Contains(err.Error(), test.wantErr) {
				t.Errorf("validate() error = %q, want it to contain %q", err, test.wantErr)
			}
		})
	}
}

func TestProfileValidateAcceptsValidEndpointConfiguration(t *testing.T) {
	profile := Profile{
		AnotherAPI: EndpointConfiguration{
			BaseURL: "https://example.com/api",
			Timeout: 1000,
		},
	}

	if err := profile.validate(); err != nil {
		t.Fatalf("validate() error = %v", err)
	}
}
