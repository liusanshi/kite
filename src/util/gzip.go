package util

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"time"
)

type streamExchange func(reader io.Reader, write io.Writer) (int64, error)

//压缩转换器
type CompressConverter struct {
	size int64
	buff *bytes.Buffer
}

// 解压转换器
type UnCompressConverter struct {
	wr        io.Writer
	size      int64
	srcReader io.Reader
}

var _ io.ReadCloser = (*CompressConverter)(nil)
var _ io.WriteCloser = (*UnCompressConverter)(nil)

// NewCompressConverter 创建一个压缩流转换器
func NewCompressConverter(wr io.ReadCloser) (*CompressConverter, error) {
	cc := &CompressConverter{}
	err := cc.init(wr)
	return cc, err
}

// NewUnCompressConverter 创建一个解压流转换器
func NewUnCompressConverter(wr io.Writer, size int64) *UnCompressConverter {
	return &UnCompressConverter{wr: wr, size: size}
}

func (cc *CompressConverter) Size() int64 {
	return cc.size
}

func (cc *CompressConverter) Close() error {
	return nil
}

func (cc *CompressConverter) init(wr io.ReadCloser) error {
	var zBuf bytes.Buffer
	zw := gzip.NewWriter(&zBuf)
	// initGzip(zw)
	_, err := io.Copy(zw, wr)
	if err != nil {
		return err
	}
	zw.Close()
	cc.size = int64(zBuf.Len())
	cc.buff = &zBuf
	return nil
}

func (cc *CompressConverter) Read(p []byte) (int, error) {
	return cc.buff.Read(p)
}

func (cc *UnCompressConverter) Close() error {
	return nil
}

// todo 实现可以支持解压的文件写入，向文件写入压缩流时，自动解压
func (cc *UnCompressConverter) Write(p []byte) (int, error) {
	// if cc.srcReader == nil {
	// 	var zBuf bytes.Buffer
	// 	zBuf.Write(p)
	// 	reader, err := gzip.NewReader(&zBuf)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// 	// cc.wr.Write()
	// }
	return cc.wr.Write(p)
}

// CompressStream2 将写入流进行压缩转换
func CompressStream2(write io.Writer) io.WriteCloser {
	zw := gzip.NewWriter(write)
	initGzip(zw)
	return zw
}

// UnCompressStream2 将读取流进行解压转换
func UnCompressStream2(reader io.Reader) (io.ReadCloser, error) {
	zr, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	return zr, nil
}

//压缩文件到可以入的流
func CompressFileToStream(filePath string, write io.Writer) (int64, error) {
	return _exchangeFileToStream(filePath, write, CompressStream)
}

//压缩文件
func CompressFile(filePath, distPath string) (int64, error) {
	return _exchangeFile(filePath, distPath, CompressStream)
}

// 初始化压缩信息
func initGzip(zw *gzip.Writer) {
	zw.Name = "kite compress"
	zw.Comment = "by payneliu"
	zw.ModTime = time.Now()
}

//使用gzip压缩流
func CompressStream(reader io.Reader, write io.Writer) (int64, error) {
	zw := gzip.NewWriter(write)
	initGzip(zw)
	defer zw.Close()
	return io.Copy(zw, reader)
}

//压缩二进制内容
func Compress(data []byte) (*bytes.Buffer, error) {
	return _exchangeData(data, CompressStream)
}

func UnCompressFileToStream(filePath string, write io.Writer) (int64, error) {
	return _exchangeFileToStream(filePath, write, UnCompressStream)
}

func UnCompressFile(filePath, distPath string) (int64, error) {
	return _exchangeFile(filePath, distPath, UnCompressStream)
}

func UnCompressStream(reader io.Reader, write io.Writer) (int64, error) {
	zr, err := gzip.NewReader(reader)
	if err != nil {
		return 0, err
	}
	defer zr.Close()
	return io.Copy(write, zr)
}

func UnCompress(data []byte) (*bytes.Buffer, error) {
	return _exchangeData(data, UnCompressStream)
}

//转换文件
func _exchangeFile(filePath, distPath string, exchage streamExchange) (int64, error) {
	distFile, err := os.OpenFile(distPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer distFile.Close()
	return _exchangeFileToStream(filePath, distFile, exchage)
}

//转换文件流
func _exchangeFileToStream(filePath string, write io.Writer, exchage streamExchange) (int64, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
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
	if _, err := exchage(&buf, &disbuf); nil != err {
		return nil, err
	}
	return &disbuf, nil
}
