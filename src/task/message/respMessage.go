package message

import (
	"fmt"
	"io"
	"net/url"
	"strconv"
)

//IMessage 消息结构体
type IMessage interface {
	io.WriterTo
	Parse(req *Request) error
}

//Type 消息类型
type Type int

const (
	// SystemMessage 系统消息
	SystemMessage = Type(0)
	// BusinessMessage 业务消息
	BusinessMessage = Type(1)
)

//curMsgID 当前的消息id
var curMsgID = 0

//检查是否实现IMessage接口
var _ IMessage = (*Message)(nil)

//Message 消息对象
type Message struct {
	ID      int
	Success bool
	Type    Type
	Content string
}

//String 将数据转换为字符串
func (m *Message) String() string {
	suc := 0
	if m.Success {
		suc = 1
	}
	return fmt.Sprintf("%d,%d,%d|%s", m.ID, suc, m.Type, m.Content)
}

//Parse 读取数据
func (m *Message) Parse(req *Request) error {
	//msg?id=1&suc=1&type=1&content=xx
	id, err := strconv.Atoi(req.Get("id"))
	if err != nil {
		return err
	}
	m.ID = id
	if req.Get("suc") == "1" {
		m.Success = true
	}
	if req.Get("type") == "1" {
		m.Type = BusinessMessage
	}
	m.Content, err = url.PathUnescape(req.Get("content"))
	if err != nil {
		return err
	}
	return nil
}

//WriteTo 写入数据
func (m *Message) WriteTo(w io.Writer) (int64, error) {
	suc := 0
	if m.Success {
		suc = 1
	}
	n, err := io.WriteString(w, fmt.Sprintf("/msg?id=%d&suc=%d&type=%d&content=%s\n", m.ID, suc, m.Type, url.PathEscape(m.Content)))
	return int64(n), err
}

//NewMessage 创建一个消息
func NewMessage(suc bool, typ Type, msg string) *Message {
	curMsgID++
	return &Message{
		Success: suc,
		Type:    typ,
		ID:      curMsgID,
		Content: msg,
	}
}
