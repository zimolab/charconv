package charconv

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"os"
)

const CreateOrTrunc = os.O_RDWR | os.O_CREATE | os.O_TRUNC

// EncoderOf 获取charsetName对应Encoder
func EncoderOf(charsetName string) *encoding.Encoder {
	en := EncodingOf(charsetName)
	if en == nil {
		return nil
	}
	return en.NewEncoder()
}

// DecoderOf 获取charsetName对应的Decoder
func DecoderOf(charsetName string) *encoding.Decoder {
	en := EncodingOf(charsetName)
	if en == nil {
		return nil
	}
	return en.NewDecoder()
}

// CodecOf 获取charsetNamed对应的Encoder和Decoder
func CodecOf(charsetName string) (*encoding.Encoder, *encoding.Decoder, error) {
	en := EncodingOf(charsetName)
	if en == nil {
		return nil, nil, unsupported(charsetName)
	}
	return en.NewEncoder(), en.NewDecoder(), nil
}

// Encode 编码字符串，即将utf-8字符串编码为指定字符集数据
func Encode(src string, destCharset string) ([]byte, error) {
	return EncodeBytes([]byte(src), destCharset)
}

// EncodeBytes 将utf-8字节切片编码为指定字符集数据
func EncodeBytes(src []byte, destCharset string) ([]byte, error) {
	if destCharset == UTF8 {
		return src, nil
	}
	// 获取编码器
	encoder := EncoderOf(destCharset)
	fmt.Println("encoder:", encoder)
	if encoder == nil {
		return nil, unsupported(destCharset)
	}
	trans := transform.NewReader(bytes.NewReader(src), encoder)
	dest, err := io.ReadAll(trans)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

// Decode 将指定字符集数据解码为utf-8字符串
func Decode(src []byte, srcCharset string) (string, error) {
	decoded, err := DecodeBytes(src, srcCharset)
	if err != nil {
		return "", err
	}
	if decoded == nil {
		return "", nil
	}
	return string(decoded), nil
}

// DecodeBytes 将指定字符集数据解码为utf-8数据
func DecodeBytes(src []byte, srcCharset string) ([]byte, error) {
	if srcCharset == UTF8 {
		return src, nil
	}
	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return nil, unsupported(srcCharset)
	}
	trans := transform.NewReader(bytes.NewReader(src), decoder)
	dest, err := io.ReadAll(trans)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

// ToUTF8Bytes DecodeBytes()函数别名
func ToUTF8Bytes(src []byte, srcCharset string) ([]byte, error) {
	return DecodeBytes(src, srcCharset)
}

// ToUTF8String Decode()函数别名
func ToUTF8String(src []byte, srcCharset string) (string, error) {
	return Decode(src, srcCharset)
}

// Convert 转换字符集
func Convert(src []byte, srcCharset string, destCharset string) (dest []byte, err error) {
	if srcCharset == destCharset {
		return src, err
	}

	if srcCharset == UTF8 {
		return Encode(string(src), destCharset)
	}

	if destCharset == UTF8 {
		decoded, err := Decode(src, srcCharset)
		if err != nil {
			return nil, err
		}
		return []byte(decoded), nil
	}

	// srcCharset != destCharset != UTF8
	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return nil, unsupported(srcCharset)
	}

	encoder := EncoderOf(destCharset)
	if encoder == nil {
		return nil, unsupported(destCharset)
	}
	// 现将src解码为utf-8这一中间形式
	srcReader := transform.NewReader(bytes.NewReader(src), decoder)
	// 再将中间形式编码为目标字符集数据
	destReader := transform.NewReader(srcReader, encoder)
	dest, err = io.ReadAll(destReader)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

// CharsetAToCharsetB Convert()函数别名
func CharsetAToCharsetB(src []byte, charsetA, charsetB string) (dest []byte, err error) {
	return Convert(src, charsetA, charsetB)
}

// EncodeTo 将utf-8数据源src编码为指定字符集形式，并写入dest中
func EncodeTo(src io.Reader, dest io.Writer, destCharset string) error {
	if destCharset == UTF8 {
		return unsupportedConversion(UTF8, destCharset)
	}
	encoder := EncoderOf(destCharset)
	if encoder == nil {
		return unsupported(destCharset)
	}
	reader := transform.NewReader(src, encoder)
	_, err := io.Copy(dest, reader)
	return err
}

func EncodeToFile(srcFilePath string, destFilePath string, destFileFlag int, destCharset string) error {
	if destCharset == UTF8 {
		return unsupportedConversion(UTF8, destCharset)
	}
	// 打开srcFile
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer closeQuietly(srcFile)
	// 打开临时文件
	tmpFile, err := os.CreateTemp("", "*")
	if err != nil {
		return err
	}
	defer removeQuietly(tmpFile)

	// 先将编码后内容写入临时文件
	err = EncodeTo(srcFile, tmpFile, destCharset)
	if err != nil {
		return err
	}
	// 再将临时文件拷贝到目标文件
	return copyTmpToDest(tmpFile, destFilePath, destFileFlag)
}

func DecodeToFile(srcFilePath string, destFilePath string, destFileFlag int, srcCharset string) error {
	if srcCharset == UTF8 {
		return unsupportedConversion(srcCharset, UTF8)
	}
	// 打开srcFile
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer closeQuietly(srcFile)
	// 打开临时文件
	tmpFile, err := os.CreateTemp("", "*")
	if err != nil {
		return err
	}
	defer removeQuietly(tmpFile)

	// 先将编码后内容写入临时文件
	err = DecodeTo(srcFile, tmpFile, srcCharset)
	if err != nil {
		return err
	}
	// 再将临时文件拷贝到目标文件
	return copyTmpToDest(tmpFile, destFilePath, destFileFlag)
}

// DecodeTo 将指定字符集数据源src编码为utf-8，并写入dest中
func DecodeTo(src io.Reader, dest io.Writer, srcCharset string) error {
	if srcCharset == UTF8 {
		return unsupportedConversion(srcCharset, UTF8)
	}
	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return unsupported(srcCharset)
	}
	reader := transform.NewReader(src, decoder)
	_, err := io.Copy(dest, reader)
	return err
}

// ConvertTo 将指定字符集数据源src转换为另一字符集，并写入dest中
func ConvertTo(src io.Reader, dest io.Writer, srcCharset, destCharset string) error {
	if srcCharset == destCharset {
		return unsupportedConversion(srcCharset, destCharset)
	}

	if srcCharset == UTF8 {
		return EncodeTo(src, dest, destCharset)
	}

	if destCharset == UTF8 {
		return DecodeTo(src, dest, srcCharset)
	}

	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return unsupported(srcCharset)
	}

	encoder := EncoderOf(destCharset)
	if encoder == nil {
		return unsupported(destCharset)
	}

	decodeReader := transform.NewReader(src, decoder)
	encodeReader := transform.NewReader(decodeReader, encoder)
	_, err := io.Copy(dest, encodeReader)
	if err != nil {
		return err
	}
	return nil
}

// ConvertFile 转换文件字符集
func ConvertFile(srcFilePath string, destFilePath string, destFileFlag int, srcCharset string, destCharset string) error {
	if srcCharset == destCharset {
		return unsupportedConversion(srcCharset, destCharset)
	}

	// utf-8 => destCharset（编码过程）
	if srcCharset == UTF8 {
		return EncodeToFile(srcFilePath, destFilePath, destFileFlag, destCharset)
	}

	// srcCharset => utf-8 (解码过程)
	if destCharset == UTF8 {
		return DecodeToFile(srcFilePath, destFilePath, destFileFlag, srcFilePath)
	}

	// srcCharset => destCharset (srcCharset、destCharset均不为utf-8，需先将srcCharset数据解码成utf-8，再重新编码为destCharset形式)
	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return unsupported(srcCharset)
	}

	encoder := EncoderOf(destCharset)
	if encoder == nil {
		return unsupported(destCharset)
	}

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer closeQuietly(srcFile)

	tmpFile, err := os.CreateTemp("", "*")
	if err != nil {
		return err
	}
	defer removeQuietly(tmpFile)

	decodeReader := transform.NewReader(srcFile, decoder)
	encodeReader := transform.NewReader(decodeReader, encoder)
	_, err = io.Copy(tmpFile, encodeReader)
	if err != nil {
		return err
	}
	return copyTmpToDest(tmpFile, destFilePath, destFileFlag)
}

func closeQuietly(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Println("error on closing file:", err)
		return
	}
}

func removeQuietly(file *os.File) {
	filename := file.Name()
	closeQuietly(file)
	err := os.Remove(filename)
	if err != nil {
		fmt.Println("error on remove file:", err)
		return
	}
}

func copyTmpToDest(tmp *os.File, destPath string, destFlag int) error {
	// 重置临时文件读取点
	_, err := tmp.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// 打开目标文件
	destFile, err := os.OpenFile(destPath, destFlag, 0666)
	if err != nil {
		return err
	}
	defer closeQuietly(destFile)

	// 将临时文件拷贝到目标文件中
	_, err = io.Copy(destFile, tmp)
	if err != nil {
		return err
	}
	return nil
}
