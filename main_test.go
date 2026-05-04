package main

import (
	"os"
	"testing"
)

func TestMainFunc(t *testing.T) {
	// Backup original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Execute main with a safe command that doesn't exit with error
	os.Args = []string{"ha-cli", "help"}
	main()
}
