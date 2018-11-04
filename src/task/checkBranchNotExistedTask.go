package task

import (
	"fmt"

	"../util"
	"./core"
)

//CheckBranchNotExistedTask 检查分支是否不存在
type CheckBranchNotExistedTask struct {
}

//检查是否实现ITask接口
var _ core.ITask = (*CheckBranchNotExistedTask)(nil)

func init() {
	util.RegisterType((*CheckBranchNotExistedTask)(nil))
}

//Init 数据初始化
func (c *CheckBranchNotExistedTask) Init(data map[string]interface{}) error {
	return nil
}

//ToMap 数据转换为map
func (c *CheckBranchNotExistedTask) ToMap() map[string]interface{} {
	return make(map[string]interface{})
}

//Run 检查分支是否不存在
func (c *CheckBranchNotExistedTask) Run(session *core.Session) error {
	_, ok := session.GetCurBranchEntity()
	if ok {
		return fmt.Errorf("branch:%s already exists", session.Branch)
	}
	return nil
}
