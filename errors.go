package charconv

import "fmt"

type ErrUnsupportedCharset struct {
	charset string
}

type ErrUnsupportedConversion struct {
	srcCharset  string
	destCharset string
}

func unsupported(charset string) ErrUnsupportedCharset {
	return ErrUnsupportedCharset{
		charset: charset,
	}
}

func unsupportedConversion(srcCharset, destCharset string) ErrUnsupportedConversion {
	return ErrUnsupportedConversion{
		srcCharset:  srcCharset,
		destCharset: destCharset,
	}
}

func (e ErrUnsupportedCharset) Error() string {
	return fmt.Sprintf("unsupported charset: %s", e.charset)
}

func (e ErrUnsupportedConversion) Error() string {
	return fmt.Sprintf("unsupported conversion: %s => %s", e.srcCharset, e.destCharset)
}
