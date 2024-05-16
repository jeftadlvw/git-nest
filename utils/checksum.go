package utils

import (
	"crypto/sha256"
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"io"
	"os"
)

/*
CalculateChecksumS calculates the checksum of a given string.
*/
func CalculateChecksumS(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

/*
CalculateChecksumF calculates the checksum of a given file.
*/
func CalculateChecksumF(path models.Path) (string, error) {
	if !path.IsFile() {
		return "", fmt.Errorf("%s is not a file", path.String())
	}

	f, err := os.Open(path.String())
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
