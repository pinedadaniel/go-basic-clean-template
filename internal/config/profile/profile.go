package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pinedadaniel/go-basic-clean-template/pkg/env"
)

type Reader func(dir string, filename string) ([]byte, error)

type OptionalConfigurations struct {
	Discriminators map[string]string `json:"discriminators,omitempty"`
}

type EndpointConfiguration struct {
	BaseURL                string                 `json:"base_url"`
	Timeout                int                    `json:"timeout"`
	CircuitBreakerRatio    float64                `json:"circuit_breaker_ratio"`
	DefaultHeaders         map[string]string      `json:"default_headers,omitempty"`
	OptionalConfigurations OptionalConfigurations `json:"optional_configurations"`
}

type Profile struct {
	AnotherAPI EndpointConfiguration `json:"another_api"`
}

func ReaderProfile(dir string, fileName string) ([]byte, error) {

	filePath := filepath.Join(dir, fileName+".json")

	absPath, err := filepath.Abs(filePath)
	if err == nil {
		filePath = absPath
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read profile file at %s: %w", filePath, err)
	}

	return bytes, nil
}

func New(reader Reader, scope env.Scope) (Profile, error) {

	if isValid, err := scope.IsScopeValid(); !isValid {
		return Profile{}, err
	}

	dir, err := DefaultConfigProfileDir()
	if err != nil {
		return Profile{}, err
	}

	raw, err := reader(dir, string(scope))
	if err != nil {
		return Profile{}, err
	}

	var profile Profile
	if err = json.Unmarshal(raw, &profile); err != nil {
		return Profile{}, fmt.Errorf("could not unmarshal profile: %w", err)
	}

	if err = profile.validate(); err != nil {
		return Profile{}, fmt.Errorf("invalid profile: %w", err)
	}

	return profile, nil
}

func DefaultConfigProfileDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get working directory %w", err)
	}
	dir = filepath.Join(dir, "internal", "config", "profile")

	return dir, nil
}

func (c *Profile) validate() error {
	var errs []error

	required := []struct {
		name           string
		timeout        int
		baseUrl        string
		circuitBreaker float64
	}{
		{"Another_api", c.AnotherAPI.Timeout, c.AnotherAPI.BaseURL, c.AnotherAPI.CircuitBreakerRatio},
	}

	for _, dep := range required {
		if dep.baseUrl == "" {
			errs = append(errs, fmt.Errorf("baseUrl not set for %s", dep.name))
		}
		if dep.timeout <= 0 {
			errs = append(errs, fmt.Errorf("timeout not set for %s", dep.name))
		}
	}

	return errors.Join(errs...)
}
