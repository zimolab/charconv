package charconv

import (
	"bytes"
	"io"
	"os"
	"testing"
)

var eucjpData = []byte{
	143, 176, 223, 185, 165, 161, 164, 192, 164, 179, 166,
}

func TestConvertWithCharsets1(t *testing.T) {
	src := bytes.NewReader(utf8Data)
	dest := MakeByteBuffer(0)
	err := ConvertBetweenCharsets(src, UTF8, dest, EUCJP)
	if err != nil {
		t.Fatal(err)
	}
	destBytes := dest.Bytes()
	if bytes.Compare(destBytes, eucjpData) != 0 {
		t.FailNow()
	}
}

func TestConvertWithCharsets2(t *testing.T) {
	src := bytes.NewReader(eucjpData)
	dest := MakeByteBuffer(0)
	err := ConvertBetweenCharsets(src, EUCJP, dest, UTF8)
	if err != nil {
		t.Fatal(err)
	}
	destBytes := dest.Bytes()
	if bytes.Compare(destBytes, utf8Data) != 0 {
		t.FailNow()
	}
}

func TestConvertWithCharsets3(t *testing.T) {
	src := bytes.NewReader(eucjpData)
	dest := MakeByteBuffer(0)
	err := ConvertBetweenCharsets(src, EUCJP, dest, GBK)
	if err != nil {
		t.Fatal(err)
	}
	destBytes := dest.Bytes()
	if bytes.Compare(destBytes, gbkData) != 0 {
		t.FailNow()
	}
}

func TestConvertWithCharsets4(t *testing.T) {
	src := bytes.NewReader(gbkData)
	dest := MakeByteBuffer(0)
	err := ConvertBetweenCharsets(src, GBK, dest, EUCJP)
	if err != nil {
		t.Fatal(err)
	}
	destBytes := dest.Bytes()
	if bytes.Compare(destBytes, eucjpData) != 0 {
		t.FailNow()
	}
}

func TestConvertFileBetweenCharsets1(t *testing.T) {
	src := "./test/test_utf8.txt"
	dest := "./test/out_eucjp.txt"
	err := ConvertFileBetweenCharsets(src, UTF8, dest, EUCJP, CreateOrTrunc)
	if err != nil {
		t.Fatal(err)
	}
	destFile, err := os.Open(dest)
	if err != nil {
		t.Fatal(err)
	}
	destBytes, err := io.ReadAll(destFile)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(destBytes, eucjpData) != 0 {
		t.Fail()
	}
}

func TestConvertFileBetweenCharsets2(t *testing.T) {
	src := "./test/out_eucjp.txt"
	dest := "./test/out_utf8_1.txt"
	err := ConvertFileBetweenCharsets(src, EUCJP, dest, UTF8, CreateOrTrunc)
	if err != nil {
		t.Fatal(err)
	}
	destFile, err := os.Open(dest)
	if err != nil {
		t.Fatal(err)
	}
	destBytes, err := io.ReadAll(destFile)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(destBytes, utf8Data) != 0 {
		t.Fail()
	}
}

func TestConvertFileBetweenCharsets3(t *testing.T) {
	src := "./test/out_eucjp.txt"
	dest := "./test/out_gbk_1.txt"
	err := ConvertFileBetweenCharsets(src, EUCJP, dest, GBK, CreateOrTrunc)
	if err != nil {
		t.Fatal(err)
	}
	destFile, err := os.Open(dest)
	if err != nil {
		t.Fatal(err)
	}
	destBytes, err := io.ReadAll(destFile)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(destBytes, gbkData) != 0 {
		t.Fail()
	}
}
