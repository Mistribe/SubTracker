package system

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type DataDownloader interface {
	Download(ctx context.Context) ([]byte, error)
	String() string
}

func newDataDownloader(source string) DataDownloader {
	if strings.HasPrefix(source, "http") ||
		strings.HasPrefix(source, "https") {
		return newHttpDownloader(source)
	}
	return newFileDownloader(source)
}

type fileDownloader struct {
	path string
}

func (f fileDownloader) String() string {
	return f.path
}

func newFileDownloader(path string) *fileDownloader {
	return &fileDownloader{path: path}
}

func (f fileDownloader) Download(ctx context.Context) ([]byte, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open the file %s: %w", f.path, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read the file %s: %w", f.path, err)
	}

	return content, nil
}

type httpDownloader struct {
	url string
}

func (h httpDownloader) String() string {
	return h.url
}

func newHttpDownloader(url string) *httpDownloader {
	return &httpDownloader{url: url}
}

func (h httpDownloader) Download(ctx context.Context) ([]byte, error) {
	resp, err := http.Get(h.url)
	if err != nil {
		return nil, fmt.Errorf("failed to download from URL %s: %w", h.url, err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", h.url, err)
	}

	return content, nil
}
