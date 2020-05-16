package tin

import (
	"fmt"
	"os"
	"strings"
)

// Config represents the configuration.
type Config struct {
	GmailCredentials string
	GmailToken       string
}

// DefaultConfig returns a tin.Config with default values.
func DefaultConfig() Config {
	home, _ := os.UserHomeDir()
	dir := fmt.Sprintf("/%v/.config/tin", strings.TrimLeft(home, "/"))

	return Config{
		GmailCredentials: dir + "/gmail/credentials.json",
		GmailToken:       dir + "/gmail/token.json",
	}
}
