package charconv

import "testing"

func TestGuessBest(t *testing.T) {
	result, err := GuessBest(utf8Data)
	if err != nil {
		t.Fatal(err)
	}
	if !charsetEquals(result.Charset, UTF8) {
		t.Fail()
	}
}

func TestGuessBestOf(t *testing.T) {
	src := "./test/test_utf8.txt"
	result, err := GuessBestOfFile(src, 1024)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
