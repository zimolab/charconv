package charconv

import (
	"bytes"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"os"
)

func Decode(src io.Reader, dest io.Writer, decoder *encoding.Decoder) error {
	srcReader := transform.NewReader(src, decoder)
	_, err := io.Copy(dest, srcReader)
	if err != nil {
		return err
	}
	return nil
}

func DecodeWithCharset(src io.Reader, dest io.Writer, srcCharset string) error {
	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return unsupported(srcCharset)
	}
	return Decode(src, dest, decoder)
}

func DecodeBytes(src []byte, dest io.Writer, decoder *encoding.Decoder) error {
	return Decode(bytes.NewReader(src), dest, decoder)
}

func DecodeBytesWithCharset(src []byte, dest io.Writer, srcCharset string) error {
	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return unsupported(srcCharset)
	}
	return DecodeWithCharset(bytes.NewReader(src), dest, srcCharset)
}

func DecodeToBytes(src io.Reader, initBuffSize int, decoder *encoding.Decoder) ([]byte, error) {
	buffer := MakeByteBuffer(initBuffSize)
	err := Decode(src, buffer, decoder)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func DecodeToBytesWithCharset(src io.Reader, initBuffSize int, srcCharset string) ([]byte, error) {
	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return nil, unsupported(srcCharset)
	}
	return DecodeToBytes(src, initBuffSize, decoder)
}

func DecodeBytesToBytes(src []byte, initBuffSize int, decoder *encoding.Decoder) ([]byte, error) {
	return DecodeToBytes(bytes.NewReader(src), initBuffSize, decoder)
}

func DecodeBytesToBytesWithCharset(src []byte, initBuffSize int, srcCharset string) ([]byte, error) {
	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return nil, unsupported(srcCharset)
	}
	return DecodeBytesToBytes(src, initBuffSize, decoder)
}

func DecodeFile(srcFilePath string, dest io.Writer, decoder *encoding.Decoder) error {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer CloseQuietly(srcFile)
	return Decode(srcFile, dest, decoder)
}

func DecodeFileWithCharset(srcFilePath string, dest io.Writer, srcCharset string) error {
	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return unsupported(srcCharset)
	}
	return DecodeFile(srcFilePath, dest, decoder)
}

func DecodeFileToFile(srcFilePath string, destFilePath string, destFileFlag int, decoder *encoding.Decoder) error {
	tmpFile, err := MakeTempFile()
	if err != nil {
		return err
	}
	defer RemoveQuietly(tmpFile)
	err = DecodeFile(srcFilePath, tmpFile, decoder)
	if err != nil {
		return err
	}
	err = tmpFile.Sync()
	if err != nil {
		return err
	}
	return CopyTmpFileTo(tmpFile, destFilePath, destFileFlag)
}

func DecodeFileToFileWithCharset(srcFilePath string, destFilePath string, destFileFlag int, srcCharset string) error {
	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return unsupported(srcCharset)
	}
	return DecodeFileToFile(srcFilePath, destFilePath, destFileFlag, decoder)
}

func DecodeToFile(src io.Reader, destFilePath string, destFileFlag int, decoder *encoding.Decoder) error {
	tmpFile, err := MakeTempFile()
	if err != nil {
		return err
	}
	defer RemoveQuietly(tmpFile)
	err = Decode(src, tmpFile, decoder)
	if err != nil {
		return err
	}
	err = tmpFile.Sync()
	if err != nil {
		return err
	}
	return CopyTmpFileTo(tmpFile, destFilePath, destFileFlag)
}

func DecodeToFileWithCharset(src io.Reader, destFilePath string, destFileFlag int, srcCharset string) error {
	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return unsupported(srcCharset)
	}
	return DecodeToFile(src, destFilePath, destFileFlag, decoder)

}

func DecodeBytesToFile(src []byte, destFilePath string, destFileFlag int, decoder *encoding.Decoder) error {
	return DecodeToFile(bytes.NewReader(src), destFilePath, destFileFlag, decoder)
}

func DecodeBytesToFileWithCharset(src []byte, destFilePath string, destFileFlag int, srcCharset string) error {
	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return unsupported(srcCharset)
	}
	return DecodeBytesToFile(src, destFilePath, destFileFlag, decoder)
}
