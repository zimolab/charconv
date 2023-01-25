package charconv

import (
	"errors"
	"github.com/saintfish/chardet"
	"io"
	"os"
)

// GuessBest 猜测data编码，返回最为接近的结果
func GuessBest(data []byte) (result *chardet.Result, err error) {
	detector := chardet.NewTextDetector()
	best, err := detector.DetectBest(data)
	if err != nil {
		return nil, err
	}
	return best, nil
}

// GuessBestOf 通过前bytesToDetect个字节，猜测src编码，返回最接近的结果
func GuessBestOf(src io.Reader, bytesToDetect int) (result *chardet.Result, err error) {
	buffer := make([]byte, bytesToDetect)
	_, err = io.ReadFull(src, buffer)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, err
	}
	detector := chardet.NewTextDetector()
	best, err := detector.DetectBest(buffer)
	if err != nil {
		return nil, err
	}
	return best, nil
}

// GuessBestOfFile 通过前bytesToDetect个字节猜测文件编码，返回最接近的结果
func GuessBestOfFile(filePath string, bytesToDetect int) (*chardet.Result, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer closeQuietly(file)
	return GuessBestOf(file, bytesToDetect)
}
