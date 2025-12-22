package docling

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL, client: &http.Client{
		Timeout: 30 * time.Second,
	}}
}

func (c *Client) Fetch(ctx context.Context, url string) (string, error) {
	// 1. Download content from URL
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	// 2. Send to Docling
	return c.Process(ctx, "downloaded.html", body)
}

func (c *Client) Process(ctx context.Context, filename string, content []byte) (string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(fw, bytes.NewReader(content)); err != nil {
		return "", err
	}
	w.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/process", &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("docling failed: %d", res.StatusCode)
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Text, nil
}
