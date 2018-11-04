package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"../../util"
)

// Branch 分支信息
type Branch struct {
	// Name 分支名称
	Name string `json:"name"`
	// Path 路径
	Path string `json:"path"`
	// Version 版本号
	Version int `json:"version"`
	// Time 最后更新时间
	Time string `json:"time"`
}

// BranchManager 分支管理
type BranchManager struct {
	// List 分支信息
	list []*Branch
	// Path 保存路径
	Path string
	// 自旋锁
	util.SpinLock
}

// Save 保存配置
func (c *BranchManager) Save() error {
	data, err := json.Marshal(c.list)
	if err != nil {
		log.Fatalf("配置文件序列化失败：%v", err)
		return err
	}
	err = ioutil.WriteFile(c.Path, data, os.ModePerm)
	if err != nil {
		log.Fatalf("配置文件保存失败：%v", err)
	}
	return err
}

// Load 加载配置
func (c *BranchManager) Load() error {
	content, err := ioutil.ReadFile(c.Path)
	if err != nil {
		log.Fatalf("配置文件:%s 加载失败:\n%v", c.Path, err)
		return err
	}
	if len(content) > 0 {
		err = json.Unmarshal(content, &c.list)
		if err != nil {
			log.Fatalf("配置文件解析失败: %v", err)
			return err
		}
	}
	return nil
}

// GetBranch 获取分支
func (c *BranchManager) GetBranch(name string) (*Branch, bool) {
	for _, item := range c.list {
		if item.Name == name {
			return item, true
		}
	}
	return nil, false
}

//DelBranch 删除分支
func (c *BranchManager) DelBranch(name string) {
	if len(name) <= 0 {
		return
	}
	list := make([]*Branch, 0, len(c.list)-1)
	for _, item := range c.list {
		if item.Name != name {
			list = append(list, item)
		}
	}
	c.list = list
}

//AddBranch 添加分支
func (c *BranchManager) AddBranch(name, path string) {
	c.list = append(c.list, &Branch{
		Name:    name,
		Version: 1,
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Path:    path,
	})
}

//Foreach 遍历所有的分支
func (c *BranchManager) Foreach(f func(*Branch, int) bool) {
	for i, item := range c.list {
		if !f(item, i) {
			break
		}
	}
}

//Filter 过滤分支
func (c *BranchManager) Filter(f func(*Branch, int) bool) []*Branch {
	list := []*Branch{}
	for i, item := range c.list {
		if f(item, i) {
			list = append(list, item)
		}
	}
	return list
}

//NewBranchManager 创建一个分支管理者
func NewBranchManager(path string) *BranchManager {
	ctx := &BranchManager{Path: path}
	err := ctx.Load()
	if err != nil {
		return nil
	}
	return ctx
}
