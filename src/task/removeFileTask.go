package task

import (
	"os"
	"path/filepath"

	"../util"
	"./core"
)

//RemoveFileTask 删除文件任务
type RemoveFileTask struct {
	// Path string
}

//检查是否实现ITask接口
var _ core.ITask = (*RemoveFileTask)(nil)

func init() {
	util.RegisterType((*RemoveFileTask)(nil))
}

//Init 数据初始化
func (r *RemoveFileTask) Init(data map[string]interface{}) error {
	// var ok bool
	// if r.Path, ok = data["Path"].(string); !ok {
	// 	return fmt.Errorf("RemoveFileTask Path type error")
	// }
	return nil
}

//ToMap 数据转换为map
func (r *RemoveFileTask) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	// data["Path"] = r.Path
	return data
}

//Run 删除文件
func (r *RemoveFileTask) Run(session *core.Session) error {
	return os.RemoveAll(filepath.Join(session.WorkSpace, session.Branch))
}
