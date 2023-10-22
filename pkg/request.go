package pkg

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Method    string
	Url       string
	Query     map[string][]string
	Body      string
	Headers   map[string]string
	BasicAuth []string // basic authentication: [username, password]
	Bearer    string   // Bearer authentication
}

func (r *Request) Do() ([]byte, error) {
	// 创建HTTP客户端
	client := &http.Client{}
	// 创建HTTP请求
	var body io.Reader
	if r.Method == "POST" {
		body = strings.NewReader(r.Body)
	}
	// 添加query参数
	if len(r.Query) > 0 {
		query := url.Values(r.Query)
		r.Url += "?" + query.Encode()
	}
	req, err := http.NewRequest(r.Method, r.Url, body)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %s %w", r.Url, err)
	}
	// 添加Basic Authentication标头
	if len(r.BasicAuth) == 2 {
		req.SetBasicAuth(r.BasicAuth[0], r.BasicAuth[1])
	}
	// 添加Bearer Authentication标头
	if r.Bearer != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.Bearer))
	}
	// 添加其他标头
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	// 发起HTTP请求

	log.Printf("发起HTTP请求: %s %s %+v %s\n", r.Method, r.Url, r.Headers, r.Body)

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求出错: %s %w", r.Url, err)
	}
	defer response.Body.Close()
	// 解析响应
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP响应解析错误: %s %w", r.Url, err)
	}

	log.Printf("HTTP响应: %s %d %s\n", r.Url, response.StatusCode, string(responseBody))

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP响应错误: %s %d %s", r.Url, response.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

func NewRequest() *Request {
	return &Request{}
}
