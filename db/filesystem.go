package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const dbFolder = "data"

// Exists returns a boolean indicating if the database directory is already initialized
func Exists() bool {
	// Create path to MANIFEST file in database directory.
	// This MANIFEST file is good indication of whether the database is initialized
	manifest := filepath.Join(Dir(), "MANIFEST")

	// Check if the MANIFEST file exists
	if _, err := os.Stat(manifest); errors.Is(err, os.ErrNotExist) {
		// File does not exist, therefore no DB
		return false
	}

	return true
}

// Dir returns the path to the directory that contains the database contents.
// It is always in the same directory as the running binary.
func Dir() string {
	// Get path to executable
	executable, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("database directory detection failure: exec path detection failure: %w", err))
	}

	// Get directory of the executable
	execDir := filepath.Dir(executable)
	// Add dbFolder to return directory of database
	return filepath.Join(execDir, dbFolder)
}
