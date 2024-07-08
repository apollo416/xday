package internal

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"os"
)

func filebase64sha256(path string) string {
	arch, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := arch.Close(); err != nil {
			log.Fatalf("Failed to close file: %v", err)
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, arch); err != nil {
		log.Fatalf("Failed to copy file content to the hash: %v", err)
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
