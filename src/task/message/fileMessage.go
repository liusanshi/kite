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
	Length int64
	Path   string
	Branch string
	md5    string
	file   io.ReadWriteCloser
}

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
	f.file = req.file
	return nil
}

//WriteTo 写入数据
func (f *FileMessage) WriteTo(w io.Writer) (int64, error) {
	// cmd?length=xx&path=xx
	// _, err := io.WriteString(w, fmt.Sprintf("%d/%s\n", f.Length, f.Path))
	num, err := io.WriteString(w,
		fmt.Sprintf("/upload?length=%d&path=%s&branch=%s&md5=%s\n",
			f.Length,
			url.PathEscape(filepath.ToSlash(filepath.Join(f.Branch, f.Path))), //将系统路径转换"/"
			url.PathEscape(f.Branch),
			f.md5))
	return int64(num), err
	//先只是写入md5，后面md5不一致再写入文件
	// if err != nil {
	// 	return 0, err
	// }
	// return io.Copy(w, f.file)
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
	_, err = io.Copy(file, io.LimitReader(f.file, f.Length))
	if err == nil || err == io.EOF {
		// fmt.Printf("upload success:%s\n", path)
		return nil
	}
	return err
}

//NewFileMessage 文件消息
func NewFileMessage(fpath, localpath, dstPath, branch string) (*FileMessage, error) {
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
	return &FileMessage{
		Length: info.Size(),
		Path:   filepath.Join(dstPath, util.Splite(fpath, localpath)),
		Branch: branch,
		file:   file,
		md5:    md5,
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
