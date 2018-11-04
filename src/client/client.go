package client

import (
	"context"
	"fmt"
	"io"
	"os"

	_ "../task" //只加载不执行
	"../task/core"
	"../util"
)

//Client 执行命令
func Client(path, cmd, branch, work string) {
	if len(path) == 0 {
		path = util.GetCurrentPath()
	}
	cfgPath := path + "/task_client.json"
	if !util.FileExists(cfgPath) {
		fmt.Printf("配置文件:%s 不存在\n", cfgPath)
		return
	}
	taskMap := core.NewMap()
	err := core.Load(cfgPath, &taskMap)
	if err != nil {
		fmt.Printf("任务加载失败: %v\n", err)
		return
	}
	taskList, ok := taskMap[cmd]
	if !ok {
		fmt.Printf("任务:%v不存在\n", cmd)
		return
	}
	session := core.NewSession(context.Background(), "root", os.Stdout, nil)
	session.TaskName = cmd
	session.Branch = branch
	session.WorkSpace = work
	err = taskList.Run(session)
	if err != nil && err != io.EOF {
		fmt.Printf("任务执行失败: %v\n", err)
		return
	}
}
