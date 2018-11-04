package task

import (
	"fmt"
	"io"
	"log"
	"net"

	"../util"
	"./core"
	"./message"
)

//TCPServerTask tcp服务的任务
type TCPServerTask struct {
	Port     string
	TaskDict core.Map
}

//检查是否实现ITask接口
var _ core.ITask = (*TCPServerTask)(nil)

//Init 数据的初始化
func (t *TCPServerTask) Init(data map[string]interface{}) error {
	var ok bool
	if t.Port, ok = data["Port"].(string); !ok {
		return fmt.Errorf("TCPServerTask Port type error")
	}
	t.TaskDict = core.NewMap()
	if dict, ok := data["TaskDict"].(map[string]interface{}); ok {
		return t.TaskDict.Init(dict)
	}
	return fmt.Errorf("TCPServerTask Port type error")
}

//ToMap 数据转换为map
func (t *TCPServerTask) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["Port"] = t.Port
	data["TaskDict"] = t.TaskDict.ToMap()
	return data
}

//Run 监听端口号，接收请求，然后根据指令执行任务；将任务的结果输出给客户端
func (t *TCPServerTask) Run(session *core.Session) error {
	listen, err := net.Listen("tcp", ":"+t.Port)
	if err != nil {
		log.Print(err)
		return err
	}
	for {
		if session.IsCancel() {
			return core.ErrCANCEL
		}
		conn, err := listen.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go t.handleConn(session.Copy(conn) /* 复制一个新的会话 */, conn)
	}
}

//handleConn 处理请求
func (t *TCPServerTask) handleConn(session *core.Session, conn net.Conn) {
	n, err := session.Request().ParseForm(conn)
	defer conn.Close()
	if err != nil && err != io.EOF {
		log.Print(err)
		session.Printf(false, message.SystemMessage, "analysis fail:%v", err)
		return
	}
	if n <= 2 {
		session.Printf(false, message.SystemMessage, "param is empty")
		return
	}
	session.Request().SetAddr(conn.RemoteAddr())
	cmd := session.Request().Cmd()
	session.Branch = session.Request().Branch()
	if task, ok := t.TaskDict[cmd]; ok {
		err := task.Run(session)
		if err != nil {
			log.Print(err)
			session.Printf(false, message.SystemMessage, "method：%s; execute fail:%v", cmd, err)
			return
		}
		log.Printf("method：%s; execute success\n", cmd)
		//执行成功
		session.Printf(true, message.SystemMessage, "method：%s; execute success", cmd)
	} else {
		log.Printf("method：%s; not fount\n", cmd)
		session.Printf(false, message.SystemMessage, "method：%s; not fount", cmd)
	}
}

func init() {
	util.RegisterType((*TCPServerTask)(nil))
}
