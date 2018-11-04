package core

import (
	"encoding/json"
	"fmt"
	"log"

	"../../util"
)

//Map 任务字典
type Map map[string]List

// Init 根据数据初始化
func (m *Map) Init(data map[string]interface{}) error {
	*m = NewMap()
	for i, item := range data {
		if task, ok := item.([]interface{}); ok {
			nt := List{}
			err := nt.Init(task)
			if err != nil {
				log.Printf("Map - subTask; err:%v\n", err)
				return err
			}
			(*m)[i] = nt
		} else if _, ok := item.(string); ok { //__type__的情况忽略
			continue
		} else {
			return fmt.Errorf("Map List subTask type error")
		}
	}
	return nil
}

//ToMap 将列表数据序列化为[]interface{}
func (m *Map) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	for key, val := range *m {
		data[key] = val.ToArray()
	}
	return data
}

//MarshalJSON 序列化
func (m *Map) MarshalJSON() ([]byte, error) {
	data := m.ToMap()
	return json.Marshal(data)
}

//UnmarshalJSON 反序列化
func (m *Map) UnmarshalJSON(data []byte) error {
	dict := make(map[string]interface{})
	err := json.Unmarshal(data, &dict)
	if err != nil {
		return err
	}
	return m.Init(dict)
}

//Run 任务运行
func (m *Map) Run(session *Session) error {
	for _, val := range *m {
		if session.IsCancel() {
			return ErrCANCEL
		}
		if err := val.Run(session); err != nil {
			return err
		}
	}
	return nil
}

// NewMap 创建一个Map
func NewMap() Map {
	return Map(make(map[string]List))
}

func init() {
	util.RegisterType((*Map)(nil))
}
