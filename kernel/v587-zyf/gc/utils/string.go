package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/v587-zyf/gc/enums"
)

// MD5 md5加密
func MD5(src string) string {
	w := md5.New()
	w.Write([]byte(src))
	return hex.EncodeToString(w.Sum(nil))
}

// GUID 产生新的GUID
func GUID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	strMD5 := base64.URLEncoding.EncodeToString(b)
	return MD5(strMD5)
}

// Token 产生新的用户登录验证码
func Token() string {
	return fmt.Sprint(time.Now().Unix(), ":", GUID())
}

func Int32ArrayToString(src []int32, flag string) (out string) {
	if len(src) == 0 {
		return ""
	}
	out = ""
	for k, v := range src {
		if k == len(src)-1 {
			out = fmt.Sprint(out, v)
		} else {
			out = fmt.Sprint(out, v, flag)
		}
	}

	return
}

func StringToIntArray(src, flag string) (out []int) {
	if src == "" {
		return nil
	}

	strs := strings.Split(src, flag)
	for _, v := range strs {
		data, err := strconv.Atoi(v)
		if err != nil {
			return nil
		}

		out = append(out, data)
	}

	return
}

func StringToFloat32Array(src, flag string) (out []float32) {
	if src == "" {
		return nil
	}

	strs := strings.Split(src, flag)
	for _, v := range strs {
		data, err := strconv.ParseFloat(v, 32)
		if err != nil {
			fmt.Printf("Failed to parse '%s' to float32: %v\n", v, err)
			return nil
		}

		out = append(out, float32(data))
	}

	return
}

func StringToInt32Array(src, flag string) (out []int32) {
	if src == "" {
		return nil
	}

	strs := strings.Split(src, flag)
	for _, v := range strs {
		data, err := strconv.Atoi(v)
		if err != nil {
			return nil
		}

		out = append(out, int32(data))
	}

	return
}

func StringToUint32Array(src, flag string) (out []uint32) {
	if src == "" {
		return nil
	}

	strs := strings.Split(src, flag)
	for _, v := range strs {
		data, err := strconv.Atoi(v)
		if err != nil {
			return nil
		}

		out = append(out, uint32(data))
	}

	return
}

func StringArrayToInt32Array(src []string) (out []int32) {
	for _, v := range src {
		data, err := strconv.Atoi(v)
		if err != nil {
			return nil
		}
		out = append(out, int32(data))
	}

	return
}

func StrToInt32(src string) int32 {
	data, err := strconv.Atoi(src)
	if err != nil {
		return 0
	}

	return int32(data)
}

func StrToFloat(src string) float64 {
	data, err := strconv.ParseFloat(src, 64)
	if err != nil {
		return 0
	}

	return data
}

func StrToFloat32(src string) float32 {
	data, err := strconv.ParseFloat(src, 32)
	if err != nil {
		return 0
	}

	return float32(data)
}

func StrToInt64(src string) int64 {
	data, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		return 0
	}

	return data
}

func StrToUInt64(src string) uint64 {
	data, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		return 0
	}

	return uint64(data)
}

func StrToUInt32(src string) uint32 {
	data, err := strconv.ParseInt(src, 10, 32)
	if err != nil {
		return 0
	}

	return uint32(data)
}

func StrToInt(src string) int {
	data, err := strconv.Atoi(src)
	if err != nil {
		return 0
	}

	return data
}
func StringToBytes(s string) []byte {
	dataPtr := (*byte)(unsafe.Pointer(&s))
	length := len(s)

	bytes := make([]byte, length)
	copy(bytes, (*[1 << 30]byte)(unsafe.Pointer(dataPtr))[:length:length])

	return bytes
}

// 表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

// 表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

// 删除空切片(字节)
func TrimSpace(s []byte) []byte {
	b := s[:0]
	for _, x := range s {
		if x != ' ' {
			b = append(b, x)
		}
	}
	return b
}

/**
 * 通过string获得一个[]string
 * @param str  10,1000,1000,1000
 * @param sep1 分隔符 ","
 */
func NewStringSlice(str string, sep string) []string {
	intSliceList := make([]string, 0)
	list := strings.Split(str, sep)
	for _, v := range list {
		intSliceList = append(intSliceList, v)
	}
	return intSliceList
}
func FnvString(s string) uint64 {
	var h = uint64(enums.Offset64)
	for _, b := range s {
		h *= enums.Prime64
		h ^= uint64(b)
	}
	return h
}

// float64ToUint64 converts a float64 to a uint64 safely
func Float64ToUint64(f float64) uint64 {
	if f < 0 {
		return 0
	}
	if f > float64(^uint64(0)) {
		return 0
	}
	return uint64(math.Floor(f))
}

// 检查字符是否是 emoji
func isEmoji(r rune) bool {
	return (r >= 0x1F600 && r <= 0x1F64F) || // 表情符号
		(r >= 0x1F300 && r <= 0x1F5FF) || // 杂项符号和象形文字
		(r >= 0x1F680 && r <= 0x1F6FF) || // 运输和地图符号
		(r >= 0x2600 && r <= 0x26FF) || // 杂项符号
		(r >= 0x2700 && r <= 0x27BF) || // 装饰符号
		(r >= 0x1F900 && r <= 0x1F9FF) || // 补充象形文字
		(r >= 0x1F1E6 && r <= 0x1F1FF) // 地区指示符号
}

// 检查字符串是否包含 emoji
func ContainsEmoji(s string) bool {
	for _, r := range s {
		if isEmoji(r) {
			return true
		}
	}
	return false
}

// 检查字符串是否包含转义
func HasInterpretedEscapeChars(s string) bool {
	for _, char := range s {
		switch char {
		case '\n', '\t', '\r', '\f', '\b', '\v', '\\':
			return true
		}
	}
	return false
}
