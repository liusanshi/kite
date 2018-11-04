package main

import (
	"flag"
	"log"
	"strings"

	"./client"
	"./config"
	"./server"
)

func main() {
	method := flag.String("func", "", `方法名称与路径
	client: 启动客户端
	server: 启动服务端`)
	fpath := flag.String("path", "", "配置文件路径")
	cmd := flag.String("cmd", "", `指令包含：
	list: 获取分支列表
	update: 更新分支
	init: 创建分支
	delete: 删除分支
	unlock: 锁住测试环境`)
	branch := flag.String("b", "", "分支名称")
	work := flag.String("workspace", "", "工作区")
	args := flag.String("args", "", "参数")

	flag.Parse()

	switch *method {
	case "config-list":
		config.List(*args, *fpath)
	case "config-set":
		params := strings.Split(*args, " ")
		if len(params) < 2 {
			log.Fatalf("参数格式错误")
		}
		config.Set(params[0], params[1], "")
	case "client":
		client.Client(*fpath, *cmd, *branch, *work)
	case "server":
		server.Sev(*fpath, *work)
	default:
		log.Println("方法名称错误")
	}
}
