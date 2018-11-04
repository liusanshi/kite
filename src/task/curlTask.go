package task

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"../util"
	"./core"
	"./message"
)

// CurlMethod curl的类型
type CurlMethod int

const (
	//GET get请求
	GET = CurlMethod(iota)
	//POST post请求
	POST = CurlMethod(iota)
)

//CurlTask curl任务
type CurlTask struct {
	//URL 请求的地址
	URL string
	//Param 请求的参数
	Param string
	//Method 请求的方法
	Method CurlMethod
	//Head http头
	Head map[string]string
}

//检查是否实现ITask接口
var _ core.ITask = (*CurlTask)(nil)

//Init 数据初始化
func (c *CurlTask) Init(data map[string]interface{}) error {
	var ok bool
	c.URL, ok = data["Url"].(string)
	if !ok {
		return fmt.Errorf("CurlTask Url type error")
	}
	c.Param, ok = data["Param"].(string)
	if !ok {
		return fmt.Errorf("CurlTask Param type error")
	}
	method := data["Method"].(float64)
	if !ok {
		return fmt.Errorf("CurlTask CurlMethod type error")
	}
	c.Method = CurlMethod(int(method))
	if data["Head"] != nil {
		c.Head = make(map[string]string)
		if head, ok := data["Head"].(map[string]interface{}); ok {
			for k, v := range head {
				c.Head[k] = v.(string)
			}
		} else {
			return fmt.Errorf("CurlTask Head type error")
		}
	}
	return nil
}

//ToMap 数据转换为map
func (c *CurlTask) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["Url"] = c.URL
	data["Param"] = c.Param
	data["Method"] = c.Method
	data["Head"] = c.Head
	return data
}

//Run 执行任务
func (c *CurlTask) Run(session *core.Session) error {
	var method = "GET"
	var url = c.URL
	var body = c.Param
	if c.Method == POST {
		method = "POST"
	} else if len(c.Param) > 0 {
		if strings.Index(c.Param, "?") > -1 {
			url += "&" + c.Param
		} else {
			url += "?" + c.Param
		}
		body = ""
	}
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.Printf("Curl NewRequest fail:%v\n", err)
		return err
	}
	if len(c.Head) > 0 {
		for k, v := range c.Head {
			request.Header.Add(k, v)
		}
	}
	if session.IsCancel() {
		return core.ErrCANCEL
	}
	request = request.WithContext(session.Ctx)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Curl request fail:%v\n", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("get url err code: %d\n", resp.StatusCode)
		return fmt.Errorf("get url err code: %d", resp.StatusCode)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Curl result resolve fail:%v\n", err)
		return err
	}
	session.Printf(true, message.SystemMessage, "%s", result)
	// session.Write(result)
	return nil
}

func init() {
	util.RegisterType((*CurlTask)(nil))
}
