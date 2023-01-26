package charconv

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"io"
	"os"
	"strings"
)

// 中文
const (
	GBK     = "GBK"
	GB18030 = "GB18030"
	Big5    = "Big5"
)

// 日文
const (
	EUCJP     = "EUC-JP"
	ISO2022JP = "ISO-2022-JP"
	ShiftJIS  = "Shift_JIS"
)

const (
	EUCKR = "EUC-KR"
)

// Unicode
const (
	UTF8    = "UTF-8"
	UTF16   = "UTF-16"
	UTF16BE = "UTF-16BE"
	UTF16LE = "UTF-16LE"
)

// 其他字符集
const (
	Macintosh = "macintosh"

	IBM037   = "IBM037"
	IBM437   = "IBM437"
	IBM850   = "IBM850"
	IBM852   = "IBM852"
	IBM855   = "IBM855"
	IBM00858 = "IBM00858"
	IBM860   = "IBM860"
	IBM862   = "IBM862"
	IBM863   = "IBM863"
	IBM865   = "IBM865"
	IBM866   = "IBM866"
	IBM1047  = "IBM1047"
	IBM01140 = "IBM01140"

	ISO88591  = "ISO-8859-1"
	ISO88592  = "ISO-8859-2"
	ISO88593  = "ISO-8859-3"
	ISO88594  = "ISO-8859-4"
	ISO88595  = "ISO-8859-5"
	ISO88596  = "ISO-8859-6"
	ISO88596E = "ISO-8859-6-E"
	ISO88596I = "ISO-8859-6-I"
	ISO88597  = "ISO-8859-7"
	ISO88598  = "ISO-8859-8"
	ISO88598E = "ISO-8859-8-E"
	ISO88598I = "ISO-8859-8-I"
	ISO88599  = "ISO-8859-9"
	ISO885910 = "ISO-8859-10"
	ISO885913 = "ISO-8859-13"
	ISO885914 = "ISO-8859-14"
	ISO885915 = "ISO-8859-15"
	ISO885916 = "ISO-8859-16"
	KOI8R     = "KOI8-R"
	KOI8U     = "KOI8-U"

	Windows874  = "Windows-874"
	Windows1250 = "Windows-1250"
	Windows1251 = "Windows-1251"
	Windows1252 = "Windows-1252"
	Windows1253 = "Windows-1253"
	Windows1254 = "Windows-1254"
	Windows1255 = "Windows-1255"
	Windows1256 = "Windows-1256"
	Windows1257 = "Windows-1257"
	Windows1258 = "Windows-1258"
)

// ianaindex.IANA.Encoding()函数无法识别GB2312、HZGB2312
// 因此需要将其映射到HZ-GB-2312。HZ-GB-2312属于GB2312的一种编码规则
var alias = map[string]string{
	"HZGB2312": "HZ-GB-2312",
	"hzgb2312": "HZ-GB-2312",
}

// EncodingOf 获取charsetName对应Encoding对象
func EncodingOf(charsetName string) encoding.Encoding {
	c, ok := alias[charsetName]
	if ok {
		charsetName = c
	}
	en, err := ianaindex.MIB.Encoding(charsetName)
	if err != nil {
		return nil
	}
	return en
}

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

// IsCharsetSupported 判断字符集是受支持
func IsCharsetSupported(charset string) bool {
	return EncodingOf(charset) != nil
}

func CopyTmpFileTo(tmp *os.File, destPath string, destFlag int) error {
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
	defer CloseQuietly(destFile)

	// 将临时文件拷贝到目标文件中
	_, err = io.Copy(destFile, tmp)
	if err != nil {
		return err
	}
	return nil
}

func CloseQuietly(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Println("error on closing file:", err)
		return
	}
}

func RemoveQuietly(file *os.File) {
	filename := file.Name()
	CloseQuietly(file)
	err := os.Remove(filename)
	if err != nil {
		fmt.Println("error on remove file:", err)
		return
	}
}

func MakeTempFile() (*os.File, error) {
	return os.CreateTemp("", "*")
}

func MakeByteBuffer(initBufferSize int) *bytes.Buffer {
	buff := make([]byte, 0, initBufferSize)
	return bytes.NewBuffer(buff)
}

func charsetEquals(charsetA, charsetB string) bool {
	return strings.ToUpper(charsetA) == strings.ToUpper(charsetB)
}
