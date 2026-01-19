package internal

import (
	"errors"
	"io"
	"os"
	"unicode/utf8"
)

// reads exactly len(data) bytes, otherwise returns an error
func ReadN(reader io.Reader, data []byte) error {
	_, err := io.ReadFull(reader, data)
	if err == io.EOF {
		return errors.New("unexpected end of file while reading")
	}
	return err
}

// writes the content to the writer
func WriteN(writer io.Writer, content []byte) error {
	_, err := writer.Write(content)
	if errors.Is(err, os.ErrClosed) {
		return errors.New("writer closed while writing")
	}
	return err
}

// checks if the encoding of the payload is valid
func CheckEncoding(enabled bool, opcode uint8, payload []byte) bool {
	if enabled && (opcode == 1 || opcode == 8) {
		return utf8.Valid(payload)
	}
	return true
}

type Buffers [][]byte

func (b Buffers) CheckEncoding(enabled bool, opcode uint8) bool {
	for i, _ := range b {
		if !CheckEncoding(enabled, opcode, b[i]) {
			return false
		}
	}
	return true
}

func (b Buffers) Len() int {
	var sum int
	for _, buffer := range b {
		sum += len(buffer)
	}
	return sum
}

// WriteTo 可重复写
func (b Buffers) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int
	for _, buffer := range b {
		n, err := w.Write(buffer)
		totalWritten += n
		if err != nil {
			return int64(totalWritten), err
		}
	}
	return int64(totalWritten), nil
}

type Bytes []byte

func (b Bytes) CheckEncoding(enabled bool, opcode uint8) bool {
	return CheckEncoding(enabled, opcode, b)
}

func (b Bytes) Len() int {
	return len(b)
}

// WriteTo 可重复写
func (b Bytes) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b)
	return int64(n), err
}
