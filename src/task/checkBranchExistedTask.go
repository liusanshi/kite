package task

import (
	"fmt"

	"../util"
	"./core"
)

//CheckBranchExistedTask 检查分支是否存在
type CheckBranchExistedTask struct {
}

//检查是否实现ITask接口
var _ core.ITask = (*CheckBranchExistedTask)(nil)

func init() {
	util.RegisterType((*CheckBranchExistedTask)(nil))
}

//Init 数据初始化
func (c *CheckBranchExistedTask) Init(data map[string]interface{}) error {
	return nil
}

//ToMap 数据转换为map
func (c *CheckBranchExistedTask) ToMap() map[string]interface{} {
	return make(map[string]interface{})
}

//Run 检查分支是否存在
func (c *CheckBranchExistedTask) Run(session *core.Session) error {
	_, ok := session.GetCurBranchEntity()
	if ok {
		return nil
	}
	return fmt.Errorf("branch:%s not exists", session.Branch)
}
