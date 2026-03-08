package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var tempFiles []string

func resolveUrl(url string) (string, error) {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(url)))
	extension := filepath.Ext(url)
	filename := fmt.Sprintf("md-slides-%s%s", hash[:16], extension)
	path := filepath.Join(os.TempDir(), filename)

	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("fetch %s: %w", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetch %s: status %d", url, response.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		return "", fmt.Errorf("write temp file: %w", err)
	}

	tempFiles = append(tempFiles, path)
	return path, nil
}

func cleanupTempFiles() {
	for _, path := range tempFiles {
		os.Remove(path)
	}
	tempFiles = nil
}
