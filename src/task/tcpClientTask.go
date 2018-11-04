package task

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"../util"
	"./core"
	"./message"
)

//TCPClientTask tcp客户端
type TCPClientTask struct {
	//IP ip
	IP string
	//Port 端口
	Port string
	//Timeout 超时时间
	Timeout int
	//Content 发送内容
	Content string
}

//检查是否实现ITask接口
var _ core.ITask = (*TCPClientTask)(nil)

//Init 数据初始化
func (t *TCPClientTask) Init(data map[string]interface{}) error {
	var ok bool
	if t.IP, ok = data["Ip"].(string); !ok {
		return fmt.Errorf("TCPClientTask Ip type error")
	}
	if t.Port, ok = data["Port"].(string); !ok {
		return fmt.Errorf("TCPClientTask Port type error")
	}
	if t.Content, ok = data["Content"].(string); !ok {
		return fmt.Errorf("TCPClientTask Content type error")
	}
	if timeout, ok := data["Timeout"].(float64); ok {
		t.Timeout = int(timeout)
	} else {
		return fmt.Errorf("TCPClientTask Content type error")
	}
	return nil
}

//ToMap 转换为map
func (t *TCPClientTask) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["Ip"] = t.IP
	data["Port"] = t.Port
	data["Content"] = t.Content
	data["Timeout"] = t.Timeout
	return data
}

//Run 执行任务
func (t *TCPClientTask) Run(session *core.Session) error {
	conn, err := net.DialTimeout("tcp", t.IP+":"+t.Port, time.Millisecond*time.Duration(t.Timeout))
	if err != nil {
		return err
	}
	defer conn.Close()
	if session.IsCancel() {
		return core.ErrCANCEL
	}
	msg := message.NewCmdMessage(t.Content, session.Branch) //创建消息
	_, err = session.Request().Send(conn, msg)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(conn)
	for {
		if session.IsCancel() {
			return core.ErrCANCEL
		}
		Msg, err := session.Response().ParseForm(reader)
		if err != nil {
			if err != io.EOF { //请求出错，且没有结束，直接返回错误
				return err
			}
			if Msg == nil || len(Msg.Content) == 0 { //读取结束，且没有数据直接返回
				return nil
			}
		}
		if !Msg.Success {
			return fmt.Errorf("%s", Msg.Content)
		}
		if Msg.Type == message.BusinessMessage {
			fmt.Fprintf(session, "%s\n", Msg.Content)
		} else {
			//todo 系统消息怎么处理? 系统消息暂时不显示
			fmt.Fprintf(os.Stderr, "%s\n", Msg.Content)
		}
	}
}

func init() {
	util.RegisterType((*TCPClientTask)(nil))
}
