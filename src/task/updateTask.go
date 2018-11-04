package task

import (
	"fmt"
	"path/filepath"
	"time"

	"../util"
	"./core"
)

//UpdateTask 更新分支的任务
type UpdateTask struct{}

//检查是否实现ITask接口
var _ core.ITask = (*UpdateTask)(nil)

func init() {
	util.RegisterType((*UpdateTask)(nil))
}

//Init 数据初始化
func (c *UpdateTask) Init(data map[string]interface{}) error {
	return nil
}

//ToMap 数据转换为map
func (c *UpdateTask) ToMap() map[string]interface{} {
	return make(map[string]interface{})
}

//Run 更新分支
func (c *UpdateTask) Run(session *core.Session) error {
	defer session.BMan.Unlock() //解锁
	b, ok := session.GetCurBranchEntity()
	if !ok {
		return fmt.Errorf("branch:%s not exist", session.Branch)
	}
	b.Version++
	b.Time = time.Now().Format("2006-01-02 15:04:05")
	b.Path = filepath.Join(session.WorkSpace, session.Branch)
	return session.BMan.Save()
}
