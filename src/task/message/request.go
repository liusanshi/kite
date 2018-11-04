package message

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/url"
	"strings"
)

//Request 请求体
type Request struct {
	cmd    string
	values url.Values
	addr   net.Addr           //客户端地址
	file   io.ReadWriteCloser //上传的文件
}

//Send 发送消息
func (r *Request) Send(w io.Writer, writeTo io.WriterTo) (int, error) {
	i, err := writeTo.WriteTo(w)
	return int(i), err
}

//SendCmd 写入数据
func (r *Request) SendCmd(w io.Writer, cmd *CmdMessage) (int, error) {
	i, err := cmd.WriteTo(w)
	return int(i), err
}

//SendFile 发送文件
func (r *Request) SendFile(w io.Writer, file *FileMessage) (int, error) {
	i, err := file.WriteTo(w)
	return int(i), err
}

//Close 关闭
func (r *Request) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

//ParseFormCmd 解析指令
func (r *Request) ParseFormCmd() (*CmdMessage, error) {
	cmd := &CmdMessage{}
	err := cmd.Parse(r)
	return cmd, err
}

//ParseFormMsg 解析消息
func (r *Request) ParseFormMsg() (*Message, error) {
	msg := &Message{}
	err := msg.Parse(r)
	return msg, err
}

//ParseFormFile 解析上传文件
func (r *Request) ParseFormFile() (*FileMessage, error) {
	file := &FileMessage{}
	err := file.Parse(r)
	return file, err
}

//ParseForm 解析请求
func (r *Request) ParseForm(read io.Reader) (int64, error) {
	nr := bufio.NewReader(read)
	head, err := nr.ReadBytes('\n')
	if err != nil {
		if err != io.EOF || len(head) == 0 {
			return 0, err
		}
	}
	head = bytes.TrimSpace(head)
	u, err := url.ParseRequestURI(string(head))
	if err != nil {
		return 0, err
	}
	r.cmd = strings.TrimLeft(u.Path, "/")
	r.values = u.Query()
	r.file = &readOnly{nr}
	return int64(len(head)), nil
}

//Cmd 获取指令
func (r *Request) Cmd() string {
	return r.cmd
}

//Branch 获取分支
func (r *Request) Branch() string {
	return r.values.Get("branch")
}

//Query 获取查询数据
func (r *Request) Query() url.Values {
	return r.values
}

//Get 获取参数
func (r *Request) Get(key string) string {
	return r.values.Get(key)
}

//Files 获取上传的文件
func (r *Request) Files() io.Reader {
	return r.file
}

//RemoteAddr 获取远程ip
func (r *Request) RemoteAddr() string {
	return strings.Split(r.addr.String(), ":")[0]
}

//Addr 获取客户端地址
func (r *Request) Addr() net.Addr {
	return r.addr
}

//SetAddr 设置远程地址
func (r *Request) SetAddr(addr net.Addr) {
	r.addr = addr
}

//NewRequest 创建一个新的请求
func NewRequest() *Request {
	return &Request{}
}

// ParseMsg 解析消息
func ParseMsg(read io.Reader) (*Message, error) {
	req := NewRequest()
	_, err := req.ParseForm(read)
	if err != nil {
		return nil, err
	}
	return req.ParseFormMsg()
}
