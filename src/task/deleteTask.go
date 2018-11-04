package task

import (
	"../util"
	"./core"
)

//DeleteTask 删除分支的任务
type DeleteTask struct{}

//检查是否实现ITask接口
var _ core.ITask = (*DeleteTask)(nil)

func init() {
	util.RegisterType((*DeleteTask)(nil))
}

//Init 数据初始化
func (c *DeleteTask) Init(data map[string]interface{}) error {
	return nil
}

//ToMap 数据转换为map
func (c *DeleteTask) ToMap() map[string]interface{} {
	return make(map[string]interface{})
}

//Run 删除分支
func (c *DeleteTask) Run(session *core.Session) error {
	defer session.BMan.Unlock() //释放锁
	session.BMan.DelBranch(session.Branch)
	return session.BMan.Save()
}
