package charconv

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"os"
)

func Convert(src io.Reader, dest io.Writer, decoder *encoding.Decoder, encoder *encoding.Encoder) error {
	decoderReader := transform.NewReader(src, decoder)
	encoderReader := transform.NewReader(decoderReader, encoder)
	_, err := io.Copy(dest, encoderReader)
	return err
}

func ConvertBetweenCharsets(src io.Reader, srcCharset string, dest io.Writer, destCharset string) error {
	if charsetEquals(srcCharset, destCharset) {
		return unsupportedConversion(srcCharset, destCharset)
	}

	if charsetEquals(srcCharset, UTF8) {
		encoder := EncoderOf(destCharset)
		if encoder == nil {
			return unsupported(destCharset)
		}
		return Encode(src, dest, encoder)
	}

	if charsetEquals(destCharset, UTF8) {
		decoder := DecoderOf(srcCharset)
		if decoder == nil {
			return unsupported(srcCharset)
		}
		return Decode(src, dest, decoder)
	}

	encoder := EncoderOf(destCharset)
	if encoder == nil {
		return unsupported(destCharset)
	}

	decoder := DecoderOf(srcCharset)
	if decoder == nil {
		return unsupported(srcCharset)
	}

	return Convert(src, dest, decoder, encoder)
}

func ConvertFileBetweenCharsets(
	srcFilePath string,
	srcFileCharset string,
	destFilePath string,
	destFileCharset string,
	destFileFlag int,
) error {
	if charsetEquals(srcFileCharset, destFileCharset) {
		return unsupportedConversion(srcFileCharset, destFileCharset)
	}

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer CloseQuietly(srcFile)

	tmpFile, err := MakeTempFile()
	if err != nil {
		return err
	}
	defer RemoveQuietly(tmpFile)

	if charsetEquals(srcFileCharset, UTF8) {
		encoder := EncoderOf(destFileCharset)
		if encoder == nil {
			return unsupported(destFileCharset)
		}
		err = Encode(srcFile, tmpFile, encoder)
		if err != nil {
			return err
		}
		err = tmpFile.Sync()
		if err != nil {
			return err
		}
		return CopyTmpFileTo(tmpFile, destFilePath, destFileFlag)

	}

	if charsetEquals(destFileCharset, UTF8) {
		decoder := DecoderOf(srcFileCharset)
		if decoder == nil {
			return unsupported(srcFileCharset)
		}
		err = Decode(srcFile, tmpFile, decoder)
		if err != nil {
			return err
		}
		err = tmpFile.Sync()
		if err != nil {
			return err
		}
		return CopyTmpFileTo(tmpFile, destFilePath, destFileFlag)
	}

	encoder := EncoderOf(destFileCharset)
	if encoder == nil {
		return unsupported(destFileCharset)
	}

	decoder := DecoderOf(srcFileCharset)
	if decoder == nil {
		return unsupported(srcFileCharset)
	}

	err = Convert(srcFile, tmpFile, decoder, encoder)
	if err != nil {
		return err
	}
	err = tmpFile.Sync()
	if err != nil {
		return err
	}
	return CopyTmpFileTo(tmpFile, destFilePath, destFileFlag)

}
