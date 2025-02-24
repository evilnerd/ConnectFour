package backend

import (
	"os"
	"path"
)

func JwtFileName() string {
	home, _ := os.UserHomeDir()
	return path.Join(home, ".connectfour-jwt")
}
