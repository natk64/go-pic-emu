// Package binary contains functions to interact with different executable formats.
package binary

import (
	"encoding/hex"
	"os"
	"regexp"
	"strings"
)

var commentRe = regexp.MustCompile(`(#.*)?\n?`)

// ParseHex parses a custom hex format file
func ParseHex(hexFile []byte) ([]byte, error) {
	str := string(hexFile)
	str = commentRe.ReplaceAllString(str, "")
	str = strings.NewReplacer(" ", "", "\n", "", "\r", "").Replace(str)
	return hex.DecodeString(str)
}

// ReadHexFile is a convenience function to read a file and parse it using [ParseHex]
func ReadHexFile(filename string) ([]byte, error) {
	hexData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseHex(hexData)
}
