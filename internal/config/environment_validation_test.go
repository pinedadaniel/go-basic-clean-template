package config

import (
	"testing"

	"github.com/pinedadaniel/go-basic-clean-template/internal/config/profile"
)

func TestValidateScope(t *testing.T) {
	tests := []struct {
		scope   string
		wantErr bool
	}{
		{scope: scopeLocal},
		{scope: scopeBeta},
		{scope: scopeProd},
		{scope: "", wantErr: true},
		{scope: "stage", wantErr: true},
		{scope: "../outside", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.scope, func(t *testing.T) {
			err := ValidateScope(test.scope)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateScope(%q) error = %v, wantErr %t", test.scope, err, test.wantErr)
			}
		})
	}
}

func TestNewRejectsNilEnvironment(t *testing.T) {
	readerCalled := false
	reader := func(string, string) ([]byte, error) {
		readerCalled = true
		return nil, nil
	}

	_, err := profile.New(reader, "local")
	if err == nil {
		t.Fatal("New() error = nil, want missing environment error")
	}
	if readerCalled {
		t.Fatal("profile reader was called without an environment")
	}
}
