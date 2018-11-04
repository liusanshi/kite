package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "../task" //只加载不执行
	"../task/core"
	"../util"
)

// Sev 服务入口
func Sev(path, work string) {
	if len(path) == 0 {
		path = util.GetCurrentPath()
	}
	cfgPath := path + "/task.json"
	if !util.FileExists(cfgPath) {
		fmt.Printf("配置文件:%s 不存在\n", cfgPath)
		return
	}
	taskList := core.NewList()
	err := core.Load(cfgPath, &taskList)
	if err != nil {
		fmt.Printf("任务加载失败: %v\n", err)
		return
	}
	branchMan := core.NewBranchManager(path + "/config.json")
	if branchMan == nil {
		return
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	session := core.NewSession(ctx, "root", os.Stdout, branchMan)
	session.WorkSpace = work
	//监听取消信号
	go func() {
		sign := listenSysSign()
		select {
		case <-sign:
			cancelFunc()
			fmt.Println("监听到取消信号")
			os.Exit(0)
		}
	}()
	err = taskList.Run(session)
	if err != nil {
		fmt.Printf("任务执行失败: %v\n", err)
		return
	}
}

// listenSysSign 监听系统退出命令
func listenSysSign() chan os.Signal {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return c
}
