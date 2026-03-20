package slurmclient

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrNotFound is returned when a requested resource does not exist.
	ErrNotFound = errors.New("slurm: resource not found")

	// ErrUnauthorized is returned when authentication fails.
	ErrUnauthorized = errors.New("slurm: unauthorized, check auth token")

	// ErrValidation is returned when client-side validation fails.
	ErrValidation = errors.New("slurm: validation error")
)

// SlurmError represents an error response from the Slurm REST API.
type SlurmError struct {
	StatusCode int
	Errors     []string
	Warnings   []string
}

func (e *SlurmError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "slurm: server returned %d", e.StatusCode)
	if len(e.Errors) > 0 {
		fmt.Fprintf(&b, ": %s", strings.Join(e.Errors, "; "))
	}
	return b.String()
}

// HasWarnings reports whether the response included warnings.
func (e *SlurmError) HasWarnings() bool {
	return len(e.Warnings) > 0
}

// RetryExhaustedError is returned when all retry attempts have failed.
type RetryExhaustedError struct {
	Attempts int
	Last     error
}

func (e *RetryExhaustedError) Error() string {
	return fmt.Sprintf("slurm: failed after %d retries: %v", e.Attempts, e.Last)
}

func (e *RetryExhaustedError) Unwrap() error {
	return e.Last
}

// VersionDiscoveryError is returned when automatic version detection fails.
type VersionDiscoveryError struct {
	Err error
}

func (e *VersionDiscoveryError) Error() string {
	return fmt.Sprintf("slurm: version discovery failed: %v", e.Err)
}

func (e *VersionDiscoveryError) Unwrap() error {
	return e.Err
}
