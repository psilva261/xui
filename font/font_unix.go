//go:build !plan9
package font

import (
	"os"
	"path"
)

func FontDir() string {
	return path.Join(os.Getenv("PLAN9"), "/font/luc")
}

func FontFile() string {
	return "unicode.18.font"
}
