package main

import (
	"testing"
)

type MockExtension struct{}

func (m *MockExtension) BeforeBuild() error { return nil }
func (m *MockExtension) AfterBuild() error  { return nil }

func TestLoadExtensions_LoadsExtensions(t *testing.T) {
	// This test is disabled because of issues with the Go plugin system in a test environment.
	t.Skip("Skipping test due to plugin issues")
}
