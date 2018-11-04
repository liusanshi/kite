package task

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"../util"
	"./core"
	"./message"
)

const (
	//maxUpload 最大上传线程
	maxUpload = 4
)

var filePathErr = fmt.Errorf("file path err") //文件遍历中断错误

//SendFileTask 发送文件的任务
type SendFileTask struct {
	//Path 本地路径
	Path string
	//DstPath 目标路径
	DstPath string
	//IP ip
	IP string
	//Port 端口
	Port string
	//排除文件
	Exclude []string
}

//检查是否实现ITask接口
var _ core.ITask = (*SendFileTask)(nil)

func init() {
	util.RegisterType((*SendFileTask)(nil))
}

//Init 数据初始化
func (s *SendFileTask) Init(data map[string]interface{}) error {
	var ok bool
	if s.Path, ok = data["Path"].(string); !ok {
		return fmt.Errorf("SendFileTask Path type error")
	}
	if s.DstPath, ok = data["DstPath"].(string); !ok {
		return fmt.Errorf("SendFileTask DstPath type error")
	}
	if s.IP, ok = data["IP"].(string); !ok {
		return fmt.Errorf("SendFileTask IP type error")
	}
	if s.Port, ok = data["Port"].(string); !ok {
		return fmt.Errorf("SendFileTask Port type error")
	}
	exclude, _ := data["Exclude"].(string)
	exclude = strings.TrimSpace(exclude)
	if len(exclude) == 0 {
		s.Exclude = []string{}
	} else {
		s.Exclude = strings.Split(exclude, " ")
	}
	return nil
}

//ToMap 数据转换为map
func (s *SendFileTask) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["IP"] = s.IP
	data["Port"] = s.Port
	data["Path"] = s.Path
	data["DstPath"] = s.DstPath
	data["Exclude"] = strings.Join(s.Exclude, " ")
	return data
}

//Run 执行任务
func (s *SendFileTask) Run(session *core.Session) error {
	var (
		errP = make(chan error)
		errC = make(chan error)
		done = make(chan struct{})
		err  error
	)
	if len(session.WorkSpace) > 0 { //如果有命令行里面携带了path，则优先使用命令行里面的path
		s.Path = session.WorkSpace
	}
	ctxP, cancelP := context.WithCancel(session.Ctx)
	ctxC, cancelC := context.WithCancel(session.Ctx)
	filepipe := s.consumerPath(ctxC, errC, done, session.Branch)
	s.productPath(ctxP, filepipe, errP)

	select {
	case err = <-errP:
		cancelC()
	case err = <-errC:
		cancelP()
	case <-done:
		break
	}
	fmt.Println("upload finish")
	return err
}

//路径生产者
func (s *SendFileTask) productPath(ctx context.Context, filepipe chan<- string, perr chan<- error) {
	go func() {
		err := filepath.Walk(s.Path, func(path string, f os.FileInfo, err error) error {
			if isEnd(ctx) {
				return filePathErr
			}
			if f == nil {
				return err
			}
			if f.Mode()&os.ModeSymlink == os.ModeSymlink { //过滤掉link文件
				return filepath.SkipDir
			}
			if f.IsDir() {
				return nil
			}
			//排除不需要的文件
			for _, ex := range s.Exclude {
				if strings.Index(path, ex) > -1 {
					return nil //这里已经是文件了，需要的是忽略，而不是跳过目录
				}
			}
			filepipe <- path
			return nil
		})
		close(filepipe)
		if err != nil {
			perr <- err
		}
	}()
}

//路径消费者
func (s *SendFileTask) consumerPath(ctx context.Context, cerr chan<- error, done chan struct{}, branch string) chan<- string {
	var filepipe = make(chan string, maxUpload)
	go func() {
		var wait sync.WaitGroup
		wait.Add(maxUpload)
		for i := 0; i < maxUpload; i++ {
			go func() {
				defer wait.Done()
				for file := range filepipe {
					if isEnd(ctx) {
						cerr <- filePathErr
						return
					}
					err := s.upload(file, branch)
					if err != nil {
						if operr, ok := err.(*net.OpError); ok {
							fmt.Printf("客户端上传错误:%v\n", operr)
							continue
						}
						cerr <- err
						return
					}
				}
			}()
		}
		wait.Wait()
		close(done)
	}()
	return filepipe
}

//upload 上传文件
func (s *SendFileTask) upload(file, branch string) error {
	conn, err := net.Dial("tcp", s.IP+":"+s.Port)
	if err != nil {
		return err
	}
	defer conn.Close()
	msg, err := message.NewFileMessage(file, s.Path, s.DstPath, branch)
	defer msg.Close()
	if err != nil {
		return err
	}
	_, err = msg.WriteTo(conn)
	if err != nil {
		return err
	}
	resp, err := message.NewResponse().ParseForm(conn)
	if err != nil {
		return err
	}
	if resp.Success {
		if resp.Content == "ok" {
			return nil
		}
		//需要上传文件
		_, err = msg.SendFile(conn)
		if err != nil {
			return err
		}
		resp, err = message.NewResponse().ParseForm(conn)
		if err != nil {
			return err
		}
		if resp.Success {
			fmt.Printf("end upload:%s\n", file)
		} else {
			fmt.Printf("upload:%s;err:%s\n", file, resp.Content)
			return fmt.Errorf(resp.Content)
		}
	} else {
		fmt.Printf("upload:%s;err:%s\n", file, resp.Content)
		return fmt.Errorf(resp.Content)
	}
	return nil
}

//isEnd 是否结束
func isEnd(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
