package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client HTTP客户端
type Client struct {
	client *http.Client
}

// NewClient 创建新的HTTP客户端
func NewClient(timeout time.Duration) *Client {
	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Post 发送POST请求
func (c *Client) Post(url string, body any, headers map[string]string) ([]byte, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	return respBody, nil
}

// Get 发送GET请求
func (c *Client) Get(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	return respBody, nil
}

// Delete 发送DELETE请求
func (c *Client) Delete(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	return respBody, nil
}
