package task

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"../util"
	"./core"
)

//替换器
type replacer struct {
	Partten string
	Repl    string
}

//Init 初始化
func (r *replacer) Init(data map[string]interface{}) error {
	var ok bool
	if r.Partten, ok = data["Partten"].(string); !ok {
		return fmt.Errorf("replacer Partten type error")
	}
	if r.Repl, ok = data["Repl"].(string); !ok {
		return fmt.Errorf("replacer Repl type error")
	}
	return nil
}

//ReplaceTask 替换文件内容的任务
type ReplaceTask struct {
	//文件路径
	FilePath string
	//替换器
	Replacer []replacer
	//编码格式
	Encoding string
}

//检查是否实现ITask接口
var _ core.ITask = (*ReplaceTask)(nil)

func init() {
	util.RegisterType((*ReplaceTask)(nil))
}

//Init 数据初始化
func (r *ReplaceTask) Init(data map[string]interface{}) error {
	var ok bool
	if r.FilePath, ok = data["FilePath"].(string); !ok {
		return fmt.Errorf("ReplaceTask FilePath type error")
	}
	var list []interface{}
	if list, ok = data["Replacer"].([]interface{}); !ok {
		return fmt.Errorf("ReplaceTask Replacer type error")
	}
	for _, repler := range list {
		if item, ok := repler.(map[string]interface{}); ok {
			repl := replacer{}
			err := repl.Init(item)
			if err != nil {
				return err
			}
			r.Replacer = append(r.Replacer, repl)
		} else {
			return fmt.Errorf("ReplaceTask sub replacer  type error")
		}
	}
	if r.Encoding, ok = data["Encoding"].(string); !ok {
		return fmt.Errorf("ReplaceTask Encoding type error")
	}
	return nil
}

//ToMap 数据转换为map
func (r *ReplaceTask) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["FilePath"] = r.FilePath
	var list []interface{}
	for _, repler := range r.Replacer {
		list = append(list, map[string]interface{}{
			"Partten": repler.Partten,
			"Repl":    repler.Repl,
		})
	}
	data["Replacer"] = list
	data["Encoding"] = r.Encoding
	return data
}

//Run 修改文件
func (r *ReplaceTask) Run(session *core.Session) error {
	if len(r.Replacer) <= 0 {
		return nil
	}
	content, err := ioutil.ReadFile(r.FilePath)
	if err != nil {
		return err
	}
	strContent := string(content)
	for _, repler := range r.Replacer {
		reg, err := regexp.Compile(session.ReplaceEnvVar(repler.Partten))
		if err != nil {
			return err
		}
		strContent = reg.ReplaceAllString(strContent, session.ReplaceEnvVar(repler.Repl))
	}
	return ioutil.WriteFile(r.FilePath, []byte(strContent), os.ModePerm)
}
