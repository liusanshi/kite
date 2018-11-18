package message

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"../../util"
)

//FileMessage 文件消息 用于上传文件
type FileMessage struct {
	Length   int64         //文件长度
	Path     string        //路径
	Branch   string        //分支
	md5      string        //文件的Md5
	Compress CompressType  //压缩类型
	file     io.ReadCloser //文件的资源地址
}

type CompressType int //压缩类型

const (
	NONE          CompressType = iota //不压缩不解压
	COMPRESSION                       //压缩
	UNCOMPRESSION                     //解压
)

//检查是否实现IMessage接口
var _ IMessage = (*FileMessage)(nil)

//String 将数据转换为字符串
func (f *FileMessage) String() string {
	return fmt.Sprintf("%s:%s:%d:%s", f.Branch, f.Path, f.Length, f.md5)
}

//Close 关闭资源
func (f *FileMessage) Close() error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

//Parse 读取数据
func (f *FileMessage) Parse(req *Request) error {
	var err error
	f.Length, err = strconv.ParseInt(req.Get("length"), 10, 0)
	if err != nil {
		return err
	}
	compress, err := strconv.ParseBool(req.Get("compress"))
	if err != nil {
		return err
	}
	f.Compress = NONE
	if compress {
		f.Compress = UNCOMPRESSION
	}
	f.Path, err = url.PathUnescape(req.Get("path"))
	if err != nil {
		return err
	}
	f.Path = filepath.FromSlash(f.Path) //将"/"转换系统路径
	f.Branch, err = url.PathUnescape(req.Get("branch"))
	if err != nil {
		return err
	}
	f.md5 = req.Get("md5")
	reader := io.LimitReader(req.file, f.Length)
	if f.Compress == UNCOMPRESSION {
		f.file = util.NewUnCompressConverter(reader)
	} else {
		f.file = util.NewNoneConverter(reader)
	}
	return nil
}

//WriteTo 写入数据
func (f *FileMessage) WriteTo(w io.Writer) (int64, error) {
	url := fmt.Sprintf("/upload?length=%d&path=%s&branch=%s&md5=%s&compress=%s\n",
		f.Length,
		url.PathEscape(filepath.ToSlash(filepath.Join(f.Branch, f.Path))), //将系统路径转换"/"
		url.PathEscape(f.Branch),
		f.md5,
		strconv.FormatBool(f.Compress == COMPRESSION))
	// log.Println(url)
	num, err := io.WriteString(w, url)
	return int64(num), err
}

// SendFile 发送文件
func (f *FileMessage) SendFile(w io.Writer) (int64, error) {
	return io.Copy(w, f.file)
}

//CheckMd5 检查文件的md5是否一致
func (f *FileMessage) CheckMd5(path string) bool {
	path = filepath.Join(path, f.Path)
	return util.FileExists(path) && util.Md5(path) == f.md5
}

// Save 保存消息
func (f *FileMessage) Save(path string) error {
	path = filepath.Join(path, f.Path)
	basename, _ := filepath.Split(path)
	//判断文件夹路径是否存在
	if !util.FileExists(basename) {
		if err := os.MkdirAll(basename, os.ModePerm); err != nil { //不存在则创建路径
			return err
		}
	} else if util.FileExists(path) && util.Md5(path) == f.md5 { //md5相同不需要上传
		return nil
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	err = file.Truncate(0) //清空文件
	if err != nil {
		return err
	}
	_, err = io.Copy(file, f.file)
	if err == nil || err == io.EOF {
		// fmt.Printf("upload success:%s\n", path)
		return nil
	}
	return err
}

//NewFileMessage 文件消息
func NewFileMessage(fpath, localpath, dstPath, branch string, isCompress bool) (*FileMessage, error) {
	if !util.FileExists(fpath) {
		return nil, os.ErrNotExist
	}
	file, err := os.OpenFile(fpath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	md5 := util.Md5(fpath)
	var (
		length int64
		fwr    io.ReadCloser
		ctype  CompressType = NONE
	)
	//压缩文件
	if isCompress {
		cc, err := util.NewCompressConverter(file)
		if err != nil {
			return nil, err
		}
		fwr = cc
		length = cc.Size()
		ctype = COMPRESSION
	} else {
		length = info.Size()
		fwr = file
	}
	return &FileMessage{
		Length:   length,
		Path:     filepath.Join(dstPath, util.Splite(fpath, localpath)),
		Branch:   branch,
		file:     fwr,
		md5:      md5,
		Compress: ctype, //压缩类型
	}, nil
}

//readOnly 只读的流
type readOnly struct {
	read io.Reader
}

func (r *readOnly) Read(p []byte) (n int, err error) {
	return r.read.Read(p)
}

func (r *readOnly) Write(p []byte) (n int, err error) {
	return
}
func (r *readOnly) Close() error {
	return nil
}
