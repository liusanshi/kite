package util

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"time"
)

type streamExchange func(reader io.Reader, write io.Writer) error

//压缩文件到可以入的流
func CompressFileToStream(filePath string, write io.Writer) error {
	return _exchangeFileToStream(filePath, write, CompressStream)
}

//压缩文件
func CompressFile(filePath, distPath string) error {
	return _exchangeFile(filePath, distPath, CompressStream)
}

// 初始化压缩信息
func initGzip(zw *gzip.Writer) {
	zw.Name = "kite compress"
	zw.Comment = "by payneliu"
	zw.ModTime = time.Now()
}

//使用gzip压缩流
func CompressStream(reader io.Reader, write io.Writer) error {
	zw := gzip.NewWriter(write)
	initGzip(zw)
	_, err := io.Copy(zw, reader)
	defer zw.Close()
	if err != nil {
		return err
	}
	return nil
}

//压缩二进制内容
func Compress(data []byte) (*bytes.Buffer, error) {
	return _exchangeData(data, CompressStream)
}

func UnCompressFileToStream(filePath string, write io.Writer) error {
	return _exchangeFileToStream(filePath, write, UnCompressStream)
}

func UnCompressFile(filePath, distPath string) error {
	return _exchangeFile(filePath, distPath, UnCompressStream)
}

func UnCompressStream(reader io.Reader, write io.Writer) error {
	zr, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	_, err = io.Copy(write, zr)
	return err
}

func UnCompress(data []byte) (*bytes.Buffer, error) {
	return _exchangeData(data, UnCompressStream)
}

//转换文件
func _exchangeFile(filePath, distPath string, exchage streamExchange) error {
	distFile, err := os.OpenFile(distPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer distFile.Close()
	return _exchangeFileToStream(filePath, distFile, exchage)
}

//转换文件流
func _exchangeFileToStream(filePath string, write io.Writer, exchage streamExchange) error {
	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return exchage(file, write)
}

//转换文件内容
func _exchangeData(data []byte, exchage streamExchange) (*bytes.Buffer, error) {
	var buf, disbuf bytes.Buffer
	if len(data) <= 0 {
		return &buf, nil
	}
	_, err := buf.Write(data)
	if err != nil {
		return nil, err
	}
	if err := exchage(&buf, &disbuf); nil != err {
		return nil, err
	}
	return &disbuf, nil
}
