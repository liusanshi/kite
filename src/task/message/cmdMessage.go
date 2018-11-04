package message

import (
	"fmt"
	"io"
)

//CmdMessage 请求结构体
type CmdMessage struct {
	//Cmd 指令
	Cmd string
	//Branch 分支
	Branch string
}

//检查是否实现IMessage接口
var _ IMessage = (*CmdMessage)(nil)

//String 将数据转换为字符串
func (cmd *CmdMessage) String() string {
	return fmt.Sprintf("%s %s", cmd.Cmd, cmd.Branch)
}

//Parse 读取数据
func (cmd *CmdMessage) Parse(req *Request) error {
	//list
	//lock
	//init?branch=11
	cmd.Cmd = req.Cmd()
	cmd.Branch = req.Get("branch")
	return nil
}

//WriteTo 写入数据
func (cmd *CmdMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, fmt.Sprintf("/%s?branch=%s\n", cmd.Cmd, cmd.Branch))
	return int64(n), err
}

//NewCmdMessage 创建一个消息
func NewCmdMessage(cmd, branch string) *CmdMessage {
	return &CmdMessage{
		Cmd:    cmd,
		Branch: branch,
	}
}
