package profile

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pinedadaniel/go-scaffolder-clean-template/pkg/env"
	"github.com/pinedadaniel/go-scaffolder-clean-template/pkg/json"
)

type OptionalConfigurations struct {
	Discriminators map[string]string `json:"discriminators,omitempty"`
}

type EndpointConfiguration struct {
	Version                string                 `json:"version,omitempty"`
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

func New(dir string, scope env.Scope) (Profile, error) {

	if isValid, err := scope.IsScopeValid(); !isValid {
		return Profile{}, err
	}

	raw, err := ReaderProfile(dir, strings.ToLower(string(scope)))
	if err != nil {
		return Profile{}, fmt.Errorf("could not load %s profile: %w", scope, err)
	}

	var profile Profile
	if err = json.Decode(raw, &profile, true); err != nil {
		return Profile{}, fmt.Errorf("could not unmarshal profile: %w", err)
	}

	if err = profile.validate(); err != nil {
		return Profile{}, fmt.Errorf("invalid profile: %w", err)
	}

	return profile, nil
}

func (c *Profile) validate() error {
	var errs []error

	required := []struct {
		name   string
		config EndpointConfiguration
	}{
		{name: "another_api", config: c.AnotherAPI},
	}

	for _, dep := range required {
		if err := validateDep(dep.name, dep.config); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func validateDep(name string, endpoint EndpointConfiguration) error {
	var errs []error

	if strings.TrimSpace(endpoint.BaseURL) == "" {
		errs = append(errs, fmt.Errorf("%s.base_url is required", name))
	} else {
		parsedURL, err := url.Parse(endpoint.BaseURL)
		if err != nil || !parsedURL.IsAbs() || parsedURL.Hostname() == "" ||
			(!strings.EqualFold(parsedURL.Scheme, "http") && !strings.EqualFold(parsedURL.Scheme, "https")) {
			errs = append(errs, fmt.Errorf("%s.base_url must be an absolute HTTP or HTTPS URL", name))
		}
	}
	if endpoint.Timeout <= 0 {
		errs = append(errs, fmt.Errorf("%s.timeout must be greater than zero", name))
	}

	return errors.Join(errs...)
}
