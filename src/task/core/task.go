package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"../../util"
)

var (
	// ErrCANCEL 取消操作的常量
	ErrCANCEL = errors.New("cancel task")
)

const (
	// TaskCONTEXTKEY 任务上下文的key
	TaskCONTEXTKEY = "_ctx_"
	// BranchCtxKey 分支名称的key
	BranchCtxKey = "branch"
)

//ITask 任务执行器
type ITask interface {
	Run(session *Session) error
	IConertToObject
	// IConertToArray
	IInit
}

//IConertToObject 转换
type IConertToObject interface {
	ToMap() map[string]interface{}
}

//IConertToArray 初始化数据的接口
type IConertToArray interface {
	ToArray() []interface{}
}

//IInit 初始化数据的接口
type IInit interface {
	Init(map[string]interface{}) error
}

//Task 任务
type Task struct {
	Type     string //任务类型
	Ignore   bool   //是否忽略错误
	Disabled bool   //是否禁用任务
	Task     ITask  //任务的实际对象
}

const (
	//TypeKey 类型的key
	TypeKey = "__type__"
)

//Init 初始化任务
func (t *Task) Init(data map[string]interface{}) error {
	var ok bool
	if t.Type, ok = data[TypeKey].(string); !ok {
		return fmt.Errorf("Task Type type error")
	}
	if ignore, ok := data["Ignore"]; ok { //忽略属性
		if ii, ok := ignore.(float64); ok {
			t.Ignore = ii == float64(1)
		} else {
			log.Printf("Task Ignore type error: require:(int);actual:(%T)", ignore)
			return fmt.Errorf("Task Ignore type error: require:(int);actual:(%T)", ignore)
		}
	}
	if disabled, ok := data["Disabled"]; ok { //禁用属性
		if ii, ok := disabled.(float64); ok {
			t.Disabled = ii == float64(1)
		} else {
			log.Printf("Task Disabled type error: require:(int);actual:(%T)", disabled)
			return fmt.Errorf("Task Disabled type error: require:(int);actual:(%T)", disabled)
		}
	}
	if temp, ok := util.NewStructPtr(t.Type); ok {
		t.Task, ok = temp.(ITask)
		if !ok {
			log.Printf("Task - iTask UnmarshalJSON; type err:%s, need ITask;\n", t.Type)
			return fmt.Errorf("Task - iTask UnmarshalJSON; type err:%s, need ITask;", t.Type)
		}
	} else {
		log.Printf("Task - iTask NewStructPtr; type err:%s;\n", t.Type)
		return fmt.Errorf("Task - iTask NewStructPtr; type err:%s;", t.Type)
	}
	err := t.Task.Init(data)
	if err != nil {
		log.Printf("Task - iTask UnmarshalJSON; err:%v\n", err)
		return err
	}
	return nil
}

// ToMap 数据转换为map
func (t *Task) ToMap() map[string]interface{} {
	data := t.Task.ToMap()
	data[TypeKey] = t.Type
	if t.Ignore {
		data["Ignore"] = 1
	} else {
		data["Ignore"] = 0
	}
	if t.Disabled {
		data["Disabled"] = 1
	} else {
		data["Disabled"] = 0
	}
	return data
}

// MarshalJSON 数据序列化
func (t *Task) MarshalJSON() ([]byte, error) {
	dic := t.ToMap()
	data, err := json.Marshal(dic)
	if err != nil {
		log.Printf("Task MarshalJSON err:%v\n", err)
		return nil, err
	}
	return data, nil
}

//UnmarshalJSON 数据反序列化
func (t *Task) UnmarshalJSON(data []byte) error {
	dic := make(map[string]interface{})
	err := json.Unmarshal(data, &dic)
	if err != nil {
		log.Printf("Task UnmarshalJSON err:%v\n", err)
		return err
	}
	err = t.Init(dic)
	if err != nil {
		log.Printf("Task - init; origin:%s; err:%v\n", data, err)
		return err
	}
	return nil
}

//Run 任务运行
func (t *Task) Run(session *Session) error {
	if t.Disabled { //禁用任务
		return nil
	}
	fmt.Printf("begin execute task:%s\n", t.Type)
	err := t.Task.Run(session)
	if t.Ignore { //忽略错误
		fmt.Printf("ignore err:%v\n", err)
		return nil
	}
	return err
}

//Load 加载数据
func Load(filePath string, unser json.Unmarshaler) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("load fail; path:%s; err:%v", filePath, err)
		return err
	}
	err = (unser).UnmarshalJSON(content)
	if err != nil {
		log.Fatalf("resolve fail; err:%v; orign:%s", err, content)
		return err
	}
	return nil
}

//Save 保存任务
func Save(filePath string, ser json.Marshaler) error {
	data, err := (ser).MarshalJSON()
	if err != nil {
		log.Fatalf("MarshalJSON fail; err:%v", err)
		return err
	}
	return ioutil.WriteFile(filePath, data, os.ModePerm)
}

// TaskWithMap 根据map转换为task
func TaskWithMap(data map[string]interface{}) (*Task, error) {
	nt := &Task{}
	err := nt.Init(data)
	if err != nil {
		return nil, err
	}
	return nt, nil
}

// TaskWithList 根据Array转换为list
func TaskWithList(data []interface{}) (*List, error) {
	nt := &List{}
	err := nt.Init(data)
	if err != nil {
		return nil, err
	}
	return nt, nil
}
