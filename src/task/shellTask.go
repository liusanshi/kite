package task

import (
	"bytes"
	"fmt"
	"os/exec"

	"../util"
	"./core"
	"./message"
)

//ShellTask shell任务
type ShellTask struct {
	Cmd  string
	Args []string
}

//检查是否实现ITask接口
var _ core.ITask = (*ShellTask)(nil)

//Init 数据初始化
func (s *ShellTask) Init(data map[string]interface{}) error {
	var ok bool
	if s.Cmd, ok = data["Cmd"].(string); !ok {
		return fmt.Errorf("ShellTask Cmd type error")
	}
	args, ok := data["Args"].([]interface{})
	if !ok {
		return fmt.Errorf("ShellTask Args type error")
	}
	for _, a := range args {
		s.Args = append(s.Args, a.(string))
	}
	return nil
}

//ToMap 数据转换为map
func (s *ShellTask) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["Cmd"] = s.Cmd
	data["Args"] = s.Args
	return data
}

//Run 执行任务
func (s *ShellTask) Run(session *core.Session) error {
	args := make([]string, len(s.Args)+1)
	args[0] = "-c"
	copy(args[1:], s.Args)
	for i, a := range args { //替换环境变量
		args[i] = session.ReplaceEnvVar(a)
	}
	cmd := exec.Command(s.Cmd, args...)
	if session.IsCancel() {
		return core.ErrCANCEL
	}
	var (
		out    bytes.Buffer
		errOut bytes.Buffer
	)
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("err:%v; info:%s", err, errOut.Bytes())
	}
	session.Printf(true, message.SystemMessage, "%s", out.Bytes())
	// session.Write(out.Bytes())
	return nil
}

func init() {
	util.RegisterType((*ShellTask)(nil))
}
