package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// Create a new hash from a file
func New(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
