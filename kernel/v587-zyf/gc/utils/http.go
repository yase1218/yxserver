package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func HttpPost(urls string, params map[string]string) ([]byte, error) {
	if !strings.Contains(urls, "http") {
		urls = "http://" + urls
	}

	values, err := json.Marshal(params)
	if err != nil {
		log.Error("json marshal err", zap.Error(err))
		return nil, err
	}
	resp, err := http.Post(urls,
		"application/x-www-form-urlencoded",
		strings.NewReader(string(values)))
	if err != nil {
		log.Error("post err", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("read body err", zap.Error(err))
		return nil, err
	}

	return body, nil
}

func HttpGet(urls string, params url.Values) ([]byte, error) {
	if !strings.Contains(urls, "http") {
		urls = "http://" + urls
	}

	parseURL, err := url.Parse(urls)
	if err != nil {
		log.Error("url parse err", zap.Error(err))
		return nil, err
	}
	parseURL.RawQuery = params.Encode()

	resp, err := http.Get(parseURL.String())
	if err != nil {
		log.Error("get err", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("read body err", zap.Error(err))
		return nil, err
	}

	return body, nil
}

// PostForm 发送http post请求数据为form
func PostForm(urls string, data url.Values) ([]byte, error) {
	if !strings.Contains(urls, "http") {
		urls = "http://" + urls
	}

	resp, err := http.PostForm(urls, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func PostJson(urls string, data []byte) ([]byte, error) {
	if !strings.Contains(urls, "http") {
		urls = "https://" + urls
	}

	var client = http.DefaultClient

	req, err := http.NewRequest("POST", urls, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func PostJsonDiyClient(urls string, data []byte, caAddr string) ([]byte, error) {
	if !strings.Contains(urls, "http") {
		urls = "https://" + urls
	}

	// 创建一个新的http.Client实例
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 跳过证书验证
			},
		},
	}

	req, err := http.NewRequest("POST", urls, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
