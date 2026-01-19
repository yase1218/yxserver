package utils

import (
	jsoniter "github.com/json-iterator/go"
	"os"
)

// WriteFile 将字符串写入文件
func WriteFile(data []byte, file string) error {
	// 写入文件 O_TRUNC 覆盖写入
	ofile, oerr := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if oerr != nil {
		return oerr
	}
	// 写完就关闭
	defer ofile.Close()

	_, werr := ofile.Write([]byte(data))
	if werr != nil {
		return werr
	}

	return nil
}

// AppendFile 将字符串写入文件(追加)
func AppendFile(data []byte, file string) error {
	// 写入文件 O_TRUNC 覆盖写入
	ofile, oerr := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if oerr != nil {
		return oerr
	}
	// 写完就关闭
	defer ofile.Close()

	_, werr := ofile.Write([]byte(data))
	if werr != nil {
		return werr
	}

	return nil
}

// WriteJSONFile 将结构转为json并记录到文件(文件不存在则创建)
func WriteJSONFile(v interface{}, file string) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// 将log status 写入文件 O_TRUNC 覆盖写入 避免一个json文件解析错误
	ofile, oerr := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if oerr != nil {
		return oerr
	}
	// 写完就关闭
	defer ofile.Close()

	_, werr := ofile.Write(body)
	if werr != nil {
		return werr
	}

	return nil
}
