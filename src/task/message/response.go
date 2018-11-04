package message

import (
	"io"
)

//Response 响应结构
type Response struct {
}

//Response 写入数据
func (r *Response) Write(w io.Writer, msg *Message) (int, error) {
	i, err := msg.WriteTo(w)
	return int(i), err
}

//ParseForm 设置body
func (r *Response) ParseForm(read io.Reader) (*Message, error) {
	return ParseMsg(read)
}

//NewResponse 创建一个新的请求
func NewResponse() *Response {
	return &Response{}
}
