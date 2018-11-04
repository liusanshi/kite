package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"../util"
)

//List 列出所有数据
func List(key, path string) (string, bool) {
	config, _ := NewConfig(path)
	config.Load()
	if len(key) == 0 {
		return "", false
	} else {
		return config.GetValue(key)
	}
}

//Set 保存key，val
func Set(ket, val, path string) {
	config, _ := NewConfig(path)
	config.Load()
	config.SetValue(ket, val)
	config.Store()
}

//Config 配置文件的结构体
type Config struct {
	LocalPath string
	Data      map[string]string
}

//NewConfig 新的配置
func NewConfig(path string) (*Config, error) {
	if len(path) == 0 {
		path = util.GetCurrentPath() + "/config.json"
	}
	if !util.FileExists(path) {
		return nil, fmt.Errorf("file:%s not exists", path)
	}
	return &Config{
		LocalPath: path,
		Data:      make(map[string]string),
	}, nil
}

//GetValue 获取数据
func (c *Config) GetValue(key string) (string, bool) {
	val, ok := c.Data[key]
	return val, ok
}

//SetValue 设置数据
func (c *Config) SetValue(key, val string) {
	c.Data[key] = val
}

//Store 保存配置
func (c *Config) Store() error {
	data, err := json.Marshal(c.Data)
	if err != nil {
		log.Fatalf("配置文件序列化失败：%v", err)
		return err
	}
	err = ioutil.WriteFile(c.LocalPath, data, os.ModePerm)
	if err != nil {
		log.Fatalf("配置文件保存失败：%v", err)
	}
	return err
}

// Load 加载配置
func (c *Config) Load() error {
	content, err := ioutil.ReadFile(c.LocalPath)
	if err != nil {
		log.Fatalf("配置文件:%s 加载失败:\n%v", c.LocalPath, err)
		return err
	}
	if len(content) > 0 {
		err = json.Unmarshal(content, &c.Data)
		if err != nil {
			log.Fatalf("配置文件解析失败: %v", err)
			return err
		}
	}
	return nil
}
