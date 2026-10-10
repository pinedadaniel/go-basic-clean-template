package json

import (
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	type config struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name                  string
		input                 string
		disallowUnknownFields bool
		wantErr               string
	}{
		{
			name:                  "valid object",
			input:                 `{"name":"service"}`,
			disallowUnknownFields: true,
		},
		{
			name:                  "unknown field rejected in strict mode",
			input:                 `{"name":"service","naem":"typo"}`,
			disallowUnknownFields: true,
			wantErr:               `unknown field "naem"`,
		},
		{
			name:    "unknown field accepted in permissive mode",
			input:   `{"name":"service","naem":"typo"}`,
			wantErr: "",
		},
		{
			name:    "multiple JSON values",
			input:   `{"name":"service"} {"name":"other"}`,
			wantErr: "unexpected trailing JSON value",
		},
		{
			name:    "trailing invalid data",
			input:   `{"name":"service"} invalid`,
			wantErr: "invalid character",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got config
			err := Decode([]byte(test.input), &got, test.disallowUnknownFields)
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("Decode() error = %v", err)
				}
				if got.Name != "service" {
					t.Errorf("Name = %q, want %q", got.Name, "service")
				}
				return
			}

			if err == nil {
				t.Fatal("Decode() error = nil")
			}
			if !strings.Contains(err.Error(), test.wantErr) {
				t.Errorf("Decode() error = %q, want it to contain %q", err, test.wantErr)
			}
		})
	}
}
