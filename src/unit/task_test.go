package unit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"kite/src/task"
	"kite/src/task/core"
	"kite/src/util"
)

//测试任务的序列化
func TestSaveTask(t *testing.T) {
	taskQueue := core.List{
		core.Task{
			Type: "CurlTask",
			Task: &task.CurlTask{
				URL:    "http://www.qq.com",
				Method: task.POST,
				Head:   map[string]string{"content-type": "html/text", "a": "1"},
			},
		},
		core.Task{
			Type: "ShellTask",
			Task: &task.ShellTask{
				Cmd:  "/usr/bin/bash",
				Args: []string{"echo hi", "echo hi-1"},
			},
		},
		core.Task{
			Type: "TCPServerTask",
			Task: &task.TCPServerTask{
				Port: "80",
				TaskDict: map[string]core.List{
					"list": {
						core.Task{
							Type: "CurlTask",
							Task: &task.CurlTask{
								URL:    "http://www.qq.com",
								Method: task.POST,
								Head:   map[string]string{"content-type": "html/text", "a": "1"},
							},
						},
						core.Task{
							Type: "ShellTask",
							Task: &task.ShellTask{
								Cmd:  "/usr/bin/bash",
								Args: []string{"echo hi", "echo hi-1"},
							},
						},
					},
					"modify": {
						core.Task{
							Type: "ShellTask",
							Task: &task.ShellTask{
								Cmd:  "/usr/bin/bash",
								Args: []string{"echo hi", "echo hi-1"},
							},
						},
					},
				},
			},
		},
	}
	path := util.GetCurrentPath() + "/task.json"
	if !util.FileExists(path) {
		_, err := os.Create(path)
		if err != nil {
			t.Logf("配置文件创建失败：%v", err)
			return
		}
	}
	err := core.Save("E:\\git\\learn\\GO\\kite\\task.json", &taskQueue)
	if err != nil {
		t.Log(err)
	}
}

//测试任务的反序列化
func TestLoadTask(t *testing.T) {
	// path := util.GetCurrentPath() + "/task.json"
	// if !util.FileExists(path) {
	// 	_, err := os.Create(path)
	// 	if err != nil {
	// 		t.Logf("配置文件创建失败：%v", err)
	// 		return
	// 	}
	// }
	taskQueue := core.NewList()
	err := core.Load("E:\\git\\learn\\GO\\kite\\task.json", &taskQueue)
	if err != nil {
		t.Log(err)
	}
	temp, _ := taskQueue.MarshalJSON()
	fmt.Printf("%s\n", temp)
	t.Logf("%s\n", temp)
}

func TestLoadClientTask(t *testing.T) {
	taskQueue := core.NewMap()
	err := core.Load("E:\\git\\learn\\GO\\kite\\task_client.json", &taskQueue)
	if err != nil {
		t.Log(err)
	}
	temp, _ := taskQueue.MarshalJSON()
	fmt.Printf("%s\n", temp)
	t.Logf("%s\n", temp)
}

func TestRegexp(t *testing.T) {
	reg, err := regexp.Compile("(?msU)###test1_brgin###.*###test1_end###")
	strContent := `<VirtualHost *>begin</VirtualHost>
	###test1_brgin###
	<VirtualHost *>
	SetEnv APP_ENV dev
	DocumentRoot /home/payneliu/git/test1/public/
	ServerName test1.qgame.qq.com
	ErrorLog logs/test1.qgame.qq.com-error_log
	CustomLog logs/test1.qgame.qq.com-access_log common
	<Directory /home/payneliu/git/test1/public/>
	Options FollowSymLinks
	AllowOverride All
	#Order allow,deny 
	#Allow from all
	</Directory>
	</VirtualHost>
	###test1_end###
	<VirtualHost *>end</VirtualHost>
	`
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(reg.ReplaceAllString(strContent, ""))
	if "" != strings.TrimSpace(reg.ReplaceAllString(strContent, "")) {
		t.Fail()
	}
}

func TestCompress(t *testing.T) {
	file, err := os.OpenFile("E:\\git\\kite\\src\\client\\client.go", os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Fatal(err)
		return
	}
	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
		return
	}
	length := info.Size()
	comp, _ := util.NewCompressConverter(file)
	data, err := ioutil.ReadAll(comp)
	err = ioutil.WriteFile("E:\\git\\kite\\src\\client\\client.go.zip", data, os.ModePerm) //到这里是对的
	if err != nil {
		t.Fatal(err)
		return
	}

	needUnCompData := bytes.NewReader(data)
	uncomp := util.NewUnCompressConverter(needUnCompData)
	data, err = ioutil.ReadAll(uncomp)
	err = ioutil.WriteFile("E:\\git\\kite\\src\\client\\client.1.go1", data, os.ModePerm) //到这里是对的
	if err != nil {
		t.Fatal(err)
		return
	}
	if length != int64(len(data)) {
		t.Fatal("长度不一致")
		return
	}
}
