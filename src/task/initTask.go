package task

import (
	"path/filepath"

	"../util"
	"./core"
)

//InitTask 创建分支的任务
type InitTask struct{}

//检查是否实现ITask接口
var _ core.ITask = (*InitTask)(nil)

func init() {
	util.RegisterType((*InitTask)(nil))
}

//Init 数据初始化
func (c *InitTask) Init(data map[string]interface{}) error {
	return nil
}

//ToMap 数据转换为map
func (c *InitTask) ToMap() map[string]interface{} {
	return make(map[string]interface{})
}

//Run 创建分支
func (c *InitTask) Run(session *core.Session) error {
	session.BMan.AddBranch(session.Branch, filepath.Join(session.WorkSpace, session.Branch)) //添加分支的地址
	defer session.BMan.Unlock()                                                              //解锁
	return session.BMan.Save()
}
