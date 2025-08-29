package updater

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Oleexo/config-go"
)

type DataDownloader interface {
	Download(ctx context.Context) ([]byte, error)
	String() string
}

func newDataDownloader(source string, cfg config.Configuration) DataDownloader {
	if strings.HasPrefix(source, "http") ||
		strings.HasPrefix(source, "https") {
		return newHttpDownloader(source, cfg)
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
	url   string
	token string
}

func (h httpDownloader) String() string {
	return h.url
}

func newHttpDownloader(url string, cfg config.Configuration) *httpDownloader {
	token := cfg.GetStringOrDefault("DATA_PRIVATE_TOKEN", "")
	return &httpDownloader{
		url:   url,
		token: token,
	}
}

func (h httpDownloader) Download(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for URL %s: %w", h.url, err)
	}
	if h.token != "" {
		if strings.HasPrefix(h.url, "https://api.github.com/") {
			req.Header.Add("Accept", "application/vnd.github.v3.raw")
			req.Header.Add("Authorization", "token "+h.token)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download from URL %s: %w", h.url, err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", h.url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download from URL %s: %s", h.url, string(content))
	}

	return content, nil
}
