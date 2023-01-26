package charconv

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestDecodeBytes(t *testing.T) {
	decoder := DecoderOf(GBK)
	if decoder == nil {
		t.Fatal()
	}
	destBuff := MakeByteBuffer(0)
	err := DecodeBytes(gbkData, destBuff, decoder)
	if err != nil {
		t.Fatal(err)
	}
	dest := destBuff.Bytes()
	if bytes.Compare(dest, utf8Data) != 0 {
		t.FailNow()
	}
}

func TestDecodeBytesWithCharset(t *testing.T) {
	destBuff := MakeByteBuffer(0)
	err := DecodeBytesWithCharset(gbkData, destBuff, GBK)
	if err != nil {
		t.Fatal(err)
	}
	dest := destBuff.Bytes()
	if bytes.Compare(dest, utf8Data) != 0 {
		t.FailNow()
	}
}

func TestDecodeToBytesWithCharset(t *testing.T) {
	dest, err := DecodeToBytesWithCharset(bytes.NewReader(gbkData), 0, GBK)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(dest, utf8Data) != 0 {
		t.FailNow()
	}
}

func TestDecodeBytesToBytesWithCharset(t *testing.T) {
	dest, err := DecodeBytesToBytesWithCharset(gbkData, 0, GBK)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(dest, utf8Data) != 0 {
		t.FailNow()
	}
}

func TestDecodeFileToFileWithCharset(t *testing.T) {
	src := "./test/test_gbk.txt"
	dest := "./test/out_utf8.txt"
	err := DecodeFileToFileWithCharset(src, dest, CreateOrTrunc, GBK)
	if err != nil {
		t.Fatal(err)
	}

	destFile, err := os.Open(dest)
	if err != nil {
		t.Fatal(err)
	}
	defer CloseQuietly(destFile)

	srcFile, err := os.Open(src)
	if err != nil {
		t.Fatal(err)
	}
	defer CloseQuietly(srcFile)

	destBytes, err := io.ReadAll(destFile)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(destBytes, utf8Data) != 0 {
		t.FailNow()
	}
}

func TestDecodeFileWithCharset(t *testing.T) {
	src := "./test/test_gbk.txt"
	dest := MakeByteBuffer(0)
	err := DecodeFileWithCharset(src, dest, GBK)
	if err != nil {
		t.Fatal(err)
	}
	destBytes := dest.Bytes()
	if bytes.Compare(destBytes, utf8Data) != 0 {
		t.FailNow()
	}
}

func TestDecodeBytesToFileWithCharset(t *testing.T) {
	dest := "./test/out_utf8.txt"
	err := DecodeBytesToFileWithCharset(utf8Data, dest, CreateOrTrunc, UTF8)
	if err != nil {
		t.Fatal(err)
	}

	destFile, err := os.Open(dest)
	if err != nil {
		t.Fatal(err)
	}
	defer CloseQuietly(destFile)

	destBytes, err := io.ReadAll(destFile)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(destBytes, utf8Data) != 0 {
		t.FailNow()
	}
}

func TestDecodeToFileWithCharset(t *testing.T) {
	dest := "./test/out_utf8.txt"
	err := DecodeToFileWithCharset(bytes.NewReader(utf8Data), dest, CreateOrTrunc, UTF8)
	if err != nil {
		t.Fatal(err)
	}

	destFile, err := os.Open(dest)
	if err != nil {
		t.Fatal(err)
	}
	defer CloseQuietly(destFile)

	destBytes, err := io.ReadAll(destFile)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(destBytes, utf8Data) != 0 {
		t.FailNow()
	}
}
