package env

import "testing"

func TestScopeIsCaseInsensitive(t *testing.T) {
	if !Scope("LOCAL").Is(Local) {
		t.Error(`Scope("LOCAL").Is(Local) = false, want true`)
	}
	if Scope("staging").Is(Local) {
		t.Error(`Scope("staging").Is(Local) = true, want false`)
	}
}

func TestScopeIsScopeValid(t *testing.T) {
	tests := []struct {
		scope   Scope
		wantErr bool
	}{
		{scope: Local},
		{scope: Beta},
		{scope: Prod},
		{scope: "LOCAL"},
		{scope: "Stage", wantErr: true},
		{scope: "../outside", wantErr: true},
		{scope: "", wantErr: true},
	}

	for _, test := range tests {
		t.Run(string(test.scope), func(t *testing.T) {
			valid, err := test.scope.IsScopeValid()
			if (err != nil) != test.wantErr {
				t.Fatalf("IsScopeValid() error = %v, wantErr %t", err, test.wantErr)
			}
			if valid == test.wantErr {
				t.Errorf("IsScopeValid() valid = %t, want %t", valid, !test.wantErr)
			}
		})
	}
}
