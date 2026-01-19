package utils

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func GetLocalIp() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func GetIpAddress(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(forwardedFor, ",")
		for _, part := range parts {
			ip := strings.TrimSpace(part)
			if ip != "" {
				return ip
			}
		}
	}
	ip := r.Header.Get("X-Real-Ip")
	if ip != "" {
		return ip
	}
	index := strings.LastIndex(r.RemoteAddr, ":")
	if index < 0 {
		return r.RemoteAddr
	}
	return r.RemoteAddr[:index]
}

func IsIPv6(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.To4() == nil
}

// GetAddrFromURL 根据URL获取网络连接地址
// Get the network connection address based on the URL
func GetAddrFromURL(URL *url.URL, tlsEnabled bool) string {
	port := SelectValue(URL.Port() == "", SelectValue(tlsEnabled, "443", "80"), URL.Port())
	hostname := URL.Hostname()
	if hostname == "" {
		hostname = "127.0.0.1"
	}
	if IsIPv6(hostname) {
		hostname = "[" + hostname + "]"
	}
	return hostname + ":" + port
}

func GetHttpIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}
