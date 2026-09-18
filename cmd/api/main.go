package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var ErrEmptyScope = errors.New("SCOPE environment variable is empty")

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Panic(ctx, "could not start app")
		log.Fatalf("launcher error: %v", err)
	}
}

func run(ctx context.Context) error {
	scope, err := getScope()
	if err != nil {
		return err
	}


	log.Info(ctx, "changing path to new dir",
		log.String("scope", scope),
	)

	if err := chdirCommand(scope); err != nil {
		return fmt.Errorf("change dir error on: %s - %w", scope, err)
	}

	cmd := execCommand("./app", os.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not start app: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("could not wait app: %w", err)
	}

	return nil
}

func getScope() (string, error) {
	value := os.Getenv("SCOPE")
	if value == "" {
		return "", ErrEmptyScope
	}

	return strings.Split(value, "-")[0], nil
}
