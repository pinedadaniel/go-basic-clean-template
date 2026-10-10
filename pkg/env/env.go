package env

import (
	"fmt"
	"strings"
)

type Scope string

const (
	Local Scope = "local"
	Beta  Scope = "beta"
	Prod  Scope = "prod"
)

func (s Scope) Is(scope Scope) bool {
	return strings.EqualFold(string(s), string(scope))
}

func (s Scope) IsScopeValid() (bool, error) {
	switch {
	case s.Is(Local), s.Is(Beta), s.Is(Prod):
		return true, nil
	default:
		return false, fmt.Errorf("invalid scope %q: allowed values are %q, %q, or %q", s, Local, Beta, Prod)
	}
}
