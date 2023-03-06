//go:build tools
// +build tools

package playlist

// pin external dependencies we install in ci
import (
	_ "github.com/joho/godotenv"
	_ "github.com/pressly/goose/cmd/goose"
)
