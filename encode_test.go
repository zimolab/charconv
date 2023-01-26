package charconv

import (
	"bytes"
	"io"
	"os"
	"testing"
)

var utf8String = "你好，世界"
var utf8Data = []byte(utf8String)

var gbkData = []byte{
	196, 227, 186, 195, 163, 172, 202, 192, 189, 231,
}

func TestCodecOf(t *testing.T) {
	decoder, encoder, err := CodecOf(GBK)
	if err != nil {
		t.Fatal(err)
	}
	if decoder == nil || encoder == nil {
		t.FailNow()
	}
}

func TestEncode(t *testing.T) {
	encoder := EncoderOf(GBK)
	if encoder == nil {
		t.Fatal()
	}
	dest := MakeByteBuffer(512)
	err := Encode(bytes.NewReader(utf8Data), dest, encoder)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	destBytes := dest.Bytes()
	//printTestData(t)
	//logDestData(t, destBytes)
	if bytes.Compare(destBytes, gbkData) != 0 {
		t.FailNow()
	}
}

func TestEncodeString(t *testing.T) {
	encoder := EncoderOf(GBK)
	if encoder == nil {
		t.Fatal()
	}
	dest := MakeByteBuffer(512)
	err := EncodeString(utf8String, dest, encoder)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	destBytes := dest.Bytes()
	//printTestData(t)
	//logDestData(t, destBytes)
	if bytes.Compare(gbkData, destBytes) != 0 {
		t.FailNow()
	}
}

func TestEncodeStringWithCharset(t *testing.T) {
	dest := MakeByteBuffer(512)
	err := EncodeStringWithCharset(utf8String, dest, GBK)
	if err != nil {
		t.Fatal(err)
	}
	destBytes := dest.Bytes()
	//printTestData(t)
	//logDestData(t, destBytes)
	if bytes.Compare(gbkData, destBytes) != 0 {
		t.FailNow()
	}
}

func TestEncodeBytes(t *testing.T) {
	encoder := EncoderOf(GBK)
	if encoder == nil {
		t.FailNow()
	}
	dest, err := EncodeToBytes(bytes.NewReader(utf8Data), 128, encoder)
	if err != nil {
		t.Fatal(err)
	}
	//printTestData(t)
	//logDestData(t, dest)
	if bytes.Compare(dest, gbkData) != 0 {
		t.FailNow()
	}
}

func TestEncodeStringToBytes(t *testing.T) {
	encoder := EncoderOf(GBK)
	if encoder == nil {
		t.FailNow()
	}
	dest, err := EncodeStringToBytes(utf8String, 0, encoder)
	if err != nil {
		t.Fatal(err)
	}
	//printTestData(t)
	//logDestData(t, dest)
	if bytes.Compare(dest, gbkData) != 0 {
		t.FailNow()
	}
}

func TestEncodeStringToBytesWithCharset(t *testing.T) {
	dest, err := EncodeStringToBytesWithCharset(utf8String, 0, GBK)
	if err != nil {
		t.Fatal(err)
	}
	//printTestData(t)
	//logDestData(t, dest)
	if bytes.Compare(dest, gbkData) != 0 {
		t.FailNow()
	}
}

func TestEncodeFileWithCharset(t *testing.T) {
	destBuffer := MakeByteBuffer(128)
	err := EncodeFileWithCharset("./test/test_utf8.txt", destBuffer, GBK)
	if err != nil {
		t.Fatal(err)
	}
	dest := destBuffer.Bytes()
	//printTestData(t)
	//logDestData(t, dest)
	if bytes.Compare(dest, gbkData) != 0 {
		t.FailNow()
	}
}

func TestEncodeFileToBytesWithCharset(t *testing.T) {
	dest, err := EncodeFileToBytesWithCharset("./test/test_utf8.txt", 128, GBK)
	if err != nil {
		t.Fatal(err)
	}
	//printTestData(t)
	//logDestData(t, dest)
	if bytes.Compare(dest, gbkData) != 0 {
		t.FailNow()
	}
}

func TestEncodeFileToFileWithCharset(t *testing.T) {
	err := EncodeFileToFileWithCharset("./test/test_utf8.txt", "./test/gbk_out.txt", CreateOrTrunc, GBK)
	if err != nil {
		t.Fatal(err)
	}
	destFile, err := os.Open("./test/gbk_out.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer CloseQuietly(destFile)
	dest, err := io.ReadAll(destFile)
	if err != nil {
		t.Fatal(err)
	}
	//printTestData(t)
	//logDestData(t, dest)
	if bytes.Compare(dest, gbkData) != 0 {
		t.FailNow()
	}
}

func printTestData(t *testing.T) {
	t.Log("utfString:", utf8String, "\n")
	t.Log("utfData:", utf8Data, "\n")
	t.Log("gbkData:", gbkData, "\n")
}

func logDestData(t *testing.T, dest []byte) {
	t.Log("dest:", dest, "\n")
}
