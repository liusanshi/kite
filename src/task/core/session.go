package core

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"kite/src/task/message"
)

// Session 会话
type Session struct {
	ID        string            //id
	Ctx       context.Context   //上下文管理器
	BMan      *BranchManager    //分支管理器
	WorkSpace string            //WorkSpace 工作路径
	TaskName  string            //TaskName 任务名称
	Args      []string          //参数
	Branch    string            //Branch 分支名称
	Compress bool              //是否启用压缩
	request   *message.Request  //request 请求对象
	response  *message.Response //response 响应对象
	write     io.Writer         //输出流
}

//Request 获取请求对象
func (c *Session) Request() *message.Request {
	if c.request == nil {
		c.request = message.NewRequest()
	}
	return c.request
}

//Response 获取响应对象
func (c *Session) Response() *message.Response {
	if c.response == nil {
		c.response = message.NewResponse()
	}
	return c.response
}

//GetCurBranchEntity 获取当前分支的实体
func (c *Session) GetCurBranchEntity() (*Branch, bool) {
	b := c.Branch
	if len(b) > 0 {
		return c.BMan.GetBranch(b)
	}
	return nil, false
}

// IsCancel 会话是否取消
func (c *Session) IsCancel() bool {
	select {
	case <-c.Ctx.Done():
		return true
	default:
		return false
	}
}

//Copy 复制一个子会话
func (c *Session) Copy(w io.Writer) *Session {
	s := &Session{
		ID:        c.ID + "/id",
		Ctx:       c.Ctx,
		BMan:      c.BMan,
		Compress: c.Compress,
		write:     w,
		WorkSpace: c.WorkSpace,
	}
	return s
}

//Write 实现io.Writer接口
func (c *Session) Write(p []byte) (n int, err error) {
	return c.write.Write(p)
}

//Printf 格式化输出
func (c *Session) Printf(suc bool, typ message.Type, format string, a ...interface{}) (n int, err error) {
	return c.Response().Write(c.write, message.NewMessage(suc, typ, fmt.Sprintf(format, a...)))
}

// ReplaceEnvVar 替换环境变量
func (c *Session) ReplaceEnvVar(repl string) string {
	branch := c.Branch
	branchPath := filepath.Join(c.WorkSpace, c.Branch)
	repl = strings.Replace(repl, "${branch}", branch, -1)
	return strings.Replace(repl, "${branchPath}", branchPath, -1)
}

//NewSession 创建一个会话
func NewSession(ctx context.Context, id string, w io.Writer, bm *BranchManager) *Session {
	return &Session{
		ID:    id,
		Ctx:   ctx,
		write: w,
		BMan:  bm,
	}
}
