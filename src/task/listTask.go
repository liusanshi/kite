package task

import (
	"../util"
	"./core"
	"./message"
)

//ListTask 显示所有分支的list
type ListTask struct{}

//检查是否实现ITask接口
var _ core.ITask = (*ListTask)(nil)

//Init 数据初始化
func (c *ListTask) Init(data map[string]interface{}) error {
	return nil
}

//ToMap 数据转换为map
func (c *ListTask) ToMap() map[string]interface{} {
	return make(map[string]interface{})
}

//Run 执行任务
func (c *ListTask) Run(session *core.Session) error {
	session.Printf(true, message.BusinessMessage, "%s\t%s\t%s", "名称", "版本", "时间")
	session.BMan.Foreach(func(b *core.Branch, i int) bool {
		session.Printf(true, message.BusinessMessage, "%s\t%d\t%s", b.Name, b.Version, b.Time)
		return true
	})
	return nil
}

func init() {
	util.RegisterType((*ListTask)(nil))
}
