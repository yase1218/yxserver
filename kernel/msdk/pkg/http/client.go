package http

//
//import (
//	"bytes"
//	"crypto/hmac"
//	"crypto/sha256"
//	"encoding/hex"
//	"encoding/json"
//	"fmt"
//	"io"
//	"net/http"
//	"net/url"
//	"strings"
//	"time"
//
//	"kernel/msdk/pkg/config"
//)
//
//// Client HTTP客户端结构
//type Client struct {
//	config     *config.Config
//	httpClient *http.Client
//}
//
//// NewClient 创建新的HTTP客户端
//func NewClient(cfg *config.Config) *Client {
//	client := &http.Client{
//		Timeout: time.Duration(cfg.HTTP.Timeout) * time.Second,
//	}
//
//	return &Client{
//		config:     cfg,
//		httpClient: client,
//	}
//}
//
//// Request HTTP请求结构
//type Request struct {
//	Method     string
//	URL        string
//	Headers    map[string]string
//	Params     url.Values
//	Body       interface{}
//	NeedEncode bool // 是否需要对请求体进行表单编码
//}
//
//// Response HTTP响应结构
//type Response struct {
//	StatusCode int
//	Headers    http.Header
//	Body       []byte
//}
//
//// Do 执行HTTP请求
//func (c *Client) Do(req *Request) (*Response, error) {
//	var (
//		httpReq *http.Request
//		err     error
//	)
//
//	// 处理请求参数
//	if req.Params != nil && len(req.Params) > 0 {
//		if strings.Contains(req.URL, "?") {
//			req.URL = req.URL + "&" + req.Params.Encode()
//		} else {
//			req.URL = req.URL + "?" + req.Params.Encode()
//		}
//	}
//
//	// 处理请求体
//	if req.Body != nil {
//		var bodyReader io.Reader
//		if req.NeedEncode {
//			// 表单编码
//			var formData url.Values
//			switch v := req.Body.(type) {
//			case url.Values:
//				formData = v
//			case map[string]string:
//				formData = url.Values{}
//				for key, val := range v {
//					formData.Set(key, val)
//				}
//			default:
//				return nil, fmt.Errorf("不支持的请求体类型，需要 url.Values 或 map[string]string")
//			}
//			bodyReader = strings.NewReader(formData.Encode())
//			if req.Headers == nil {
//				req.Headers = make(map[string]string)
//			}
//			req.Headers["Content-Type"] = "application/x-www-form-urlencoded"
//		} else {
//			// JSON编码
//			jsonData, err := json.Marshal(req.Body)
//			if err != nil {
//				return nil, fmt.Errorf("JSON编码失败: %w", err)
//			}
//			bodyReader = bytes.NewReader(jsonData)
//			if req.Headers == nil {
//				req.Headers = make(map[string]string)
//			}
//			req.Headers["Content-Type"] = "application/json"
//		}
//
//		httpReq, err = http.NewRequest(req.Method, req.URL, bodyReader)
//		if err != nil {
//			return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
//		}
//	} else {
//		httpReq, err = http.NewRequest(req.Method, req.URL, nil)
//		if err != nil {
//			return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
//		}
//	}
//
//	// 设置请求头
//	for key, value := range c.config.HTTP.Headers {
//		httpReq.Header.Set(key, value)
//	}
//	for key, value := range req.Headers {
//		httpReq.Header.Set(key, value)
//	}
//
//	// 发送请求
//	resp, err := c.httpClient.Do(httpReq)
//	if err != nil {
//		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
//	}
//	defer resp.Body.Close()
//
//	// 读取响应体
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return nil, fmt.Errorf("读取响应体失败: %w", err)
//	}
//
//	return &Response{
//		StatusCode: resp.StatusCode,
//		Headers:    resp.Header,
//		Body:       body,
//	}, nil
//}
//
//// Get 发送GET请求
//func (c *Client) Get(url string, params url.Values, headers map[string]string) (*Response, error) {
//	req := &Request{
//		Method:  "GET",
//		URL:     url,
//		Params:  params,
//		Headers: headers,
//	}
//	return c.Do(req)
//}
//
//// Post 发送POST请求
//func (c *Client) Post(url string, body interface{}, headers map[string]string, needEncode bool) (*Response, error) {
//	req := &Request{
//		Method:     "POST",
//		URL:        url,
//		Body:       body,
//		Headers:    headers,
//		NeedEncode: needEncode,
//	}
//	return c.Do(req)
//}
//
//// GenerateHMACSHA256 生成HMAC-SHA256签名
//func GenerateHMACSHA256(secret string, message string) string {
//	h := hmac.New(sha256.New, []byte(secret))
//	h.Write([]byte(message))
//	return hex.EncodeToString(h.Sum(nil))
//}
//
//// ParseJSONResponse 解析JSON响应
//func ParseJSONResponse(response *Response, v interface{}) error {
//	if response.StatusCode != http.StatusOK {
//		return fmt.Errorf("服务器返回非200状态码: %d", response.StatusCode)
//	}
//
//	err := json.Unmarshal(response.Body, v)
//	if err != nil {
//		return fmt.Errorf("解析JSON响应失败: %w", err)
//	}
//
//	return nil
//}
