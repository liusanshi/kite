package util

import (
	"reflect"
)

//typeRegistry 类型缓存
var typeRegistry = make(map[string]reflect.Type)

//NewStructPtr 创建指定类型的结构体的指针
func NewStructPtr(clsName string) (interface{}, bool) {
	if typ, ok := typeRegistry[clsName]; ok {
		return reflect.New(typ).Interface(), true
	}
	return nil, false
}

//NewStruct 创建指定类型的结构体
func NewStruct(clsName string) (interface{}, bool) {
	if typ, ok := typeRegistry[clsName]; ok {
		return reflect.New(typ).Elem().Interface(), true
	}
	return nil, false
}

//RegisterType 注册
func RegisterType(elem interface{}) {
	t := reflect.TypeOf(elem).Elem()
	typeRegistry[t.Name()] = t
}
