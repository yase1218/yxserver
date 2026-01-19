package taptap

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func OAuth(kid, macKey string) error {

	// clientId 参数在 TapDC 后台查看
	clientId := "rqjjpth0ly9c6gg5ow"

	// 随机数，正式上线请替换
	nonce := "8IBTHwOdqNKAWeKl7plt66=="

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	reqHost := "open.tapapis.cn"
	reqURI := "/account/profile/v1?client_id=" + clientId
	reqURL := "https://" + reqHost + reqURI

	macStr := timestamp + "\n" + nonce + "\n" + "GET" + "\n" + reqURI + "\n" + reqHost + "\n" + "443" + "\n\n"
	mac := hmacSha1(macStr, macKey)
	authorization := "MAC id=" + "\"" + kid + "\"" + "," + "ts=" + "\"" + timestamp + "\"" + "," + "nonce=" + "\"" + nonce + "\"" + "," + "mac=" + "\"" + mac + "\""

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 添加请求头
	req.Header.Add("Authorization", authorization)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(respBody))

	return nil
}

/*
HMAC-SHA1 签名
*/
func hmacSha1(valStr, keyStr string) string {
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(valStr))

	// 进行 Base64 编码
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
