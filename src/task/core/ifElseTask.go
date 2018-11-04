package core

import (
	"fmt"
	"log"

	"../../util"
)

// IConditions 条件接口
type IConditions interface {
	ITask
	GetResult() int //获取条件执行的结果
}

// IfElse 条件任务  分支任务
type IfElse struct {
	Cond     IConditions //条件
	Result   int         //条件的对比值
	Logic    string      //逻辑计算:=; >; <; >=; <=;
	Body     List        //完成之后的执行body
	ElseTask List        //与之匹配的else
}

//检查是否实现ITask接口
var _ ITask = (*IfElse)(nil)

func init() {
	util.RegisterType((*IfElse)(nil))
}

//Init 初始化任务
func (i *IfElse) Init(data map[string]interface{}) error {
	if val, ok := data["Cond"].(map[string]interface{}); ok {
		nt, err := TaskWithMap(val)
		if err != nil {
			log.Printf("IfElse - Cond; err:%v\n", err)
			return err
		}
		if i.Cond, ok = nt.Task.(IConditions); !ok {
			log.Printf("IfElse - Cond; type err: cant convert (%T) to (%T)\n", nt.Task, IConditions(nil))
		}
	}
	if val, ok := data["ElseTask"].([]interface{}); ok {
		nt, err := TaskWithList(val)
		if err != nil {
			log.Printf("IfElse - ElseTask; err:%v\n", err)
			return err
		}
		i.ElseTask = *nt
	}
	if val, ok := data["Body"].([]interface{}); ok {
		nt, err := TaskWithList(val)
		if err != nil {
			log.Printf("IfElse - Body; err:%v\n", err)
			return err
		}
		i.Body = *nt
	}
	if result, ok := data["Result"]; ok { //忽略属性
		if ii, ok := result.(float64); ok {
			i.Result = int(ii)
		} else {
			log.Printf("Task Result type error: require:(int);actual:(%T)", result)
			return fmt.Errorf("Task Result type error: require:(int);actual:(%T)", result)
		}
	}
	i.Logic, _ = data["Logic"].(string)
	return nil
}

// ToMap 数据转换为map
func (i *IfElse) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["Cond"] = i.Cond.ToMap()
	data["Body"] = i.Body.ToArray()
	data["ElseTask"] = i.ElseTask.ToArray()
	data["Result"] = i.Result
	data["Logic"] = i.Logic
	return data
}

//Run 任务运行
func (i *IfElse) Run(session *Session) error {
	err := i.Cond.Run(session)
	if err != nil {
		return err
	}
	if i.compu(i.Cond.GetResult()) {
		if i.Body != nil && len(i.Body) > 0 {
			return i.Body.Run(session)
		}
		return nil
	}
	if i.ElseTask != nil && len(i.ElseTask) > 0 {
		return i.ElseTask.Run(session)
	}
	return nil
}

// compu 结算结果 满足条件
func (i *IfElse) compu(result int) bool {
	switch i.Logic {
	case ">":
		return result > i.Result
	case "<":
		return result < i.Result
	case "==":
		return result == i.Result
	case ">=":
		return result >= i.Result
	case "<=":
		return result <= i.Result
	}
	panic("logic not in ('>', '<', '==', '>=', '<=')")
}
