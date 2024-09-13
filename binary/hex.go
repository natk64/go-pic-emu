// Package binary contains functions to interact with different executable formats.
package binary

import (
	"bytes"
	"encoding/hex"
	"os"
	"regexp"
	"strings"

	"github.com/marcinbor85/gohex"
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

// ParseIHex parses a ihex format file
func ParseIHex(hexFile []byte) ([]byte, error) {
	mem := gohex.NewMemory()
	if err := mem.ParseIntelHex(bytes.NewReader(hexFile)); err != nil {
		return nil, err
	}

	return mem.ToBinary(0, 64_000, 0), nil
}

// ReadIHexFile is a convenience function to read a file and parse it using [ParseIHex]
func ReadIHexFile(filename string) ([]byte, error) {
	hexData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseIHex(hexData)
}
