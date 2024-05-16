package tests

import (
	"crypto/sha256"
	"fmt"
	"github.com/jeftadlvw/git-nest/models"
	"github.com/jeftadlvw/git-nest/utils"
	"io"
	"os"
	"testing"
)

func TestCalculateChecksumS(t *testing.T) {
	cases := []struct {
		input string
	}{
		{""},
		{"foo"},
		{"7ßijß8m aß 8nads84654f687gspoilömxöoijasoiu8s ()(/)ASDnja s"},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestCalculateChecksumS-%d", index+1), func(t *testing.T) {
			h := sha256.New()
			h.Write([]byte(tc.input))

			expectedChecksum := fmt.Sprintf("%x", h.Sum(nil))
			checksum := utils.CalculateChecksumS(tc.input)

			if checksum != expectedChecksum {
				t.Fatalf("checksums do not match: expected: %s, got: %s", expectedChecksum, checksum)
			}
		})
	}
}

func TestCalculateChecksumF(t *testing.T) {
	cases := []struct {
		input string
	}{
		{""},
		{"foo"},
		{"7ßijß8m aß 8nads84654f687gspoilömxöoijasoiu8s ()(/)ASDnja s"},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestCalculateChecksumF-%d", index+1), func(t *testing.T) {

			tempDir := models.Path(t.TempDir())
			tempFile := tempDir.SJoin("checksum.txt")

			err := utils.WriteStrToFile(tempFile, tc.input)
			if err != nil {
				t.Fatalf("error writing temporary file: %s", err)
			}

			f, err := os.Open(tempFile.String())
			if err != nil {
				t.Fatalf("error opening temporary file: %s", err)
			}
			defer f.Close()

			h := sha256.New()
			if _, err = io.Copy(h, f); err != nil {
				t.Fatalf("error reading temporary file: %s", err)
			}

			expectedChecksum := fmt.Sprintf("%x", h.Sum(nil))
			checksumString := utils.CalculateChecksumS(tc.input)
			checksumFile, err := utils.CalculateChecksumF(tempFile)
			if err != nil {
				t.Fatalf("error calculating checksum for file: %s", err)
			}

			if checksumFile != expectedChecksum {
				t.Fatalf("checksums do not match: expected: %s, got: %s", expectedChecksum, checksumFile)
			}
			if checksumFile != checksumString {
				t.Fatalf("unequal checksums for file and string: %s != %s", checksumFile, checksumString)
			}
		})
	}
}
