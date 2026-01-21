package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http" // Import textproto for CanonicalMIMEHeaderKey if needed internally
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
// headers 参数现在使用自定义的 httpclient.Header 类型
func (c *Client) Post(url string, body any, headers Header) ([]byte, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置headers
	// 先设置 Content-Type
	req.Header.Set("Content-Type", "application/json")
	// 然后遍历自定义 Header 类型的实例，将其值复制到 *http.Request 的 Header 中
	// 由于我们的 Header 类型实现了 Set, Get, Add, Del 等方法，这里的 h 是 map[string][]string
	// 需要遍历每个键值对
	for key, values := range headers {
		// 遍历该键对应的所有值
		for _, value := range values {
			// req.Header 是 http.Header 类型，即 map[string][]string
			// 它的 Add 方法会追加值，Set 方法会覆盖值
			// 这里使用 Add 更符合我们自定义 Header 的多值特性
			req.Header.Add(key, value)
		}
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
// headers 参数现在使用自定义的 httpclient.Header 类型
func (c *Client) Get(url string, headers Header) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置headers
	// 遍历自定义 Header 类型的实例
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
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
// headers 参数现在使用自定义的 httpclient.Header 类型
func (c *Client) Delete(url string, headers Header) ([]byte, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置headers
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
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
