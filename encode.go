package charconv

import (
	"bytes"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"os"
)

// Encode 基础编码方法
func Encode(src io.Reader, dest io.Writer, destEncoder *encoding.Encoder) error {
	encodeReader := transform.NewReader(src, destEncoder)
	_, err := io.Copy(dest, encodeReader)
	return err
}

func EncodeString(src string, dest io.Writer, destEncoder *encoding.Encoder) error {
	return Encode(bytes.NewReader([]byte(src)), dest, destEncoder)
}

func EncodeStringWithCharset(src string, dest io.Writer, destCharset string) error {
	encoder := EncoderOf(destCharset)
	if encoder == nil {
		return unsupported(destCharset)
	}
	return EncodeString(src, dest, encoder)
}

func EncodeToBytes(src io.Reader, initBuffSize int, destEncoder *encoding.Encoder) ([]byte, error) {
	buffer := MakeByteBuffer(initBuffSize)
	err := Encode(src, buffer, destEncoder)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), err
}

func EncodeStringToBytes(src string, initBuffSize int, destEncoder *encoding.Encoder) ([]byte, error) {
	return EncodeToBytes(bytes.NewReader([]byte(src)), initBuffSize, destEncoder)
}

func EncodeStringToBytesWithCharset(src string, initBuffSize int, destCharset string) ([]byte, error) {
	encoder := EncoderOf(destCharset)
	if encoder == nil {
		return nil, unsupported(destCharset)
	}
	return EncodeStringToBytes(src, initBuffSize, encoder)
}

func EncodeFile(srcFilePath string, dest io.Writer, destEncoder *encoding.Encoder) error {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer CloseQuietly(srcFile)
	return Encode(srcFile, dest, destEncoder)
}

func EncodeFileWithCharset(srcFilePath string, dest io.Writer, destCharset string) error {
	encoder := EncoderOf(destCharset)
	if encoder == nil {
		return unsupported(destCharset)
	}
	return EncodeFile(srcFilePath, dest, encoder)
}

func EncodeFileToBytes(srcFilePath string, initBuffSize int, destEncoder *encoding.Encoder) ([]byte, error) {
	buffer := MakeByteBuffer(initBuffSize)
	err := EncodeFile(srcFilePath, buffer, destEncoder)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func EncodeFileToBytesWithCharset(srcFilePath string, initBuffSize int, destCharset string) ([]byte, error) {
	encoder := EncoderOf(destCharset)
	if encoder == nil {
		return nil, unsupported(destCharset)
	}
	return EncodeFileToBytes(srcFilePath, initBuffSize, encoder)
}

func EncodeFileToFile(srcFilePath string, destFilePath string, destFileFlag int, destEncoder *encoding.Encoder) error {
	tmpFile, err := os.CreateTemp("", "*")
	if err != nil {
		return err
	}
	defer RemoveQuietly(tmpFile)

	err = EncodeFile(srcFilePath, tmpFile, destEncoder)
	if err != nil {
		return err
	}

	return CopyTmpFileTo(tmpFile, destFilePath, destFileFlag)

}

func EncodeFileToFileWithCharset(srcFilePath string, destFilePath string, destFileFlag int, destCharset string) error {
	encoder := EncoderOf(destCharset)
	if encoder == nil {
		return unsupported(destCharset)
	}
	return EncodeFileToFile(srcFilePath, destFilePath, destFileFlag, encoder)
}
