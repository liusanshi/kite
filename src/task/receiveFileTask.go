package task

import (
	"fmt"
	"strings"

	"../util"
	"./core"
	"./message"
)

//ReceiveFileTask 接收文件的任务
type ReceiveFileTask struct {
	//Path 本地路径
	// Path string
	//Port 端口
	// Port string
	//IPLists ip白名单
	IPLists []string
}

//检查是否实现ITask接口
var _ core.ITask = (*ReceiveFileTask)(nil)

func init() {
	util.RegisterType((*ReceiveFileTask)(nil))
}

//Init 数据初始化
func (s *ReceiveFileTask) Init(data map[string]interface{}) error {
	var ok bool
	// if s.Path, ok = data["Path"].(string); !ok {
	// 	return fmt.Errorf("ReceiveFileTask Path type error")
	// }
	// if s.Port, ok = data["Port"].(string); !ok {
	// 	return fmt.Errorf("ReceiveFileTask Port type error")
	// }
	iplist, ok := data["IPLists"].(string)
	if !ok {
		return fmt.Errorf("ReceiveFileTask IPLists type error")
	}
	s.IPLists = strings.Split(iplist, " ")
	return nil
}

//ToMap 数据转换为map
func (s *ReceiveFileTask) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	// data["Port"] = s.Port
	// data["Path"] = s.Path
	data["IPLists"] = strings.Join(s.IPLists, " ")
	return data
}

//Run 保存上传的文件
func (s *ReceiveFileTask) Run(session *core.Session) error {
	ip := session.Request().RemoteAddr()
	if util.IndexOf(s.IPLists, ip) == -1 {
		return fmt.Errorf("ip:%s not in the white list", ip)
	}
	msg, err := session.Request().ParseFormFile()
	if err != nil {
		return err
	}
	if msg.CheckMd5(session.WorkSpace) { //文件md5一致，则退出接受文件
		session.Printf(true, message.SystemMessage, "ok")
		return nil
	}
	session.Printf(true, message.SystemMessage, "ready") //表示服务器已经准备好接受文件
	err = msg.Save(session.WorkSpace)
	//返回客户端是否成功
	if err == nil {
		session.Printf(true, message.SystemMessage, "ok")
	} else {
		session.Printf(false, message.SystemMessage, "%v", err)
	}
	return err
}
