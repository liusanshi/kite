## 风筝系统 - 自助部署测试环境

### 支持下列功能
1. list 列出所有的测试环境  
>示例：
```
./kite --func=client --path=/home/payneliu/git/kite/ --cmd=list --b=test1
其中：
func: client表示客户端; server表示服务端
path: 相关配置的存放位置(task_client.json)等地址
cmd: 指令; 包含：(list、init、update、delete、unlock)
b: 环境的名称(list、unlock命令时选填，其他必填)
```
2. init 创建一个测试环境  
>示例：
```
./kite --func=client --path=/home/payneliu/git/kite/ --cmd=init --b=test1
```
3. update 更新一个测试环境  
>示例：
```
./kite --func=client --path=/home/payneliu/git/kite/ --cmd=update --b=test1
```
4. delete 删除一个测试环境  
>示例：
```
./kite --func=client --path=/home/payneliu/git/kite/ --cmd=delete --b=test1
```
5. unlock 解锁一个测试环境  
>示例：
```
./kite --func=client --path=/home/payneliu/git/kite/ --cmd=unlock --b=test1
```
6. 启动服务器  
>示例：
```
./kite --func=server --path=/home/payneliu/git/kite/ --workspace=/home/payneliu/git/
其中：
func: client表示客户端; server表示服务端
path: 相关配置的存放位置(task.json, config.json)等地址
workspace: 工作目录
```

### 支持的任务列表
1. CheckBranchExistedTask
>作用：检查测试是否存在环境    
作用范围：服务端  
使用方法：
```
{
    "__type__": "CheckBranchExistedTask"
}
```

2. CheckBranchNotExistedTask
>作用：检查测试是否不存在环境  
作用范围：服务端  
使用方法：
```
{
    "__type__": "CheckBranchNotExistedTask"
}
```

3. CurlTask
>作用：curl请求；可以用于执行完成之后做一些hook操作  
作用范围：服务端; 客户端  
使用方法：
```
{
    "Url": "http://www.qq.com",
    "Param": "a=1&b=2",
    "Method": 0, //0:GET;1:POST
    "Head": [{
        "Host": "www.qq.com"
    }],
    "__type__": "CurlTask"
},
```

4. DeleteTask
>作用：删除测试环境  
作用范围：服务端  
使用方法：
```
{
    "__type__": "DeleteTask"
}
```

5. InitTask
>作用：创建测试环境  
作用范围：服务端  
使用方法：
```
{
    "__type__": "InitTask"
}
```

6. ListTask
>作用：列出所有的测试环境  
作用范围：服务端  
使用方法：
```
{
    "__type__": "ListTask"
}
```

7. LockTask
>作用：锁住所有环境  
作用范围：服务端  
使用方法：
```
{
    "__type__": "LockTask"
}
```

8. ReceiveFileTask
>作用：接收文件  
作用范围：服务端  
使用方法：
```
{
    "IPLists": "127.0.0.1", //ip白名单 空格分隔
    "__type__": "ReceiveFileTask"
}
```

9. RemoveFileTask
>作用：删除环境的所有文件  
作用范围：服务端  
使用方法：
```
{
    "__type__": "RemoveFileTask"
}
```

10. ReplaceTask
>作用：替换文本的任务  
作用范围：服务端、客户端  
使用方法：
```
{
    "FilePath": "/data/home/payneliu/services/apache-2.4/conf/httpd.conf", //文件地址
    "Encoding": "utf8", //编码格式
    "Replacer": [{
        "Partten": "(?msU)###${branch}_begin###.*###${branch}_end###", //匹配模式
        "Repl": ""                                                     //替换的字符串
    }],
    "__type__": "ReplaceTask"
}
```

11. SendFileTask
>作用：客户端发送文件  
作用范围：客户端  
使用方法：
```
{
    "Path": "/data/home/payneliu/git/crayfish", //本地目录
    "DstPath": "/",                             //服务器目录
    "Exclude": "/vendor/",                      //排除匹配的文件或目录，使用空格分隔多个
    "IP": "127.0.0.1",                          //上传的目标ip
    "Port": "8880",                             //上传服务的port
    "__type__": "SendFileTask"
}
```

12. ShellTask
>作用：执行shell脚本  
作用范围：服务端、客户端  
使用方法：
```
{
    "Args": ["php7 /home/payneliu/git/crayfish/artisan route:cache"], //shell命令
    "Cmd": "/bin/bash",
    "Ignore": 0,
    "__type__": "ShellTask"
}
```

13. TCPClientTask
>作用：主要的客户端任务；他负责将指令发送给服务端  
作用范围：客户端  
使用方法：
```
{
    "Ip": "127.0.0.1",  //服务端ip
    "Port": "8880",     //服务端port
    "Content": "lock",  //指令的内容
    "Timeout": 3000,    //超时时间(单位：ms)
    "__type__": "TCPClientTask"
}
```

14. TCPServerTask
>作用：主要的服务端任务，是所有服务端任务的宿主；他负责处理客服端的指令，进而解析为各个任务再执行  
作用范围：服务端  
使用方法：
```
[{
    "Port": "8880",
    "__type__": "TCPServerTask"
    "TaskDict": {  //所有的任务字典，客户端的命令，根据TaskDict找到具体的指令
        "list": [{
            "__type__": "ListTask"
        }]
    }
}]
```

15. UnlockTask
>作用：解锁整个测试环境  
作用范围：服务端  
使用方法：
```
{
    "__type__": "UnlockTask"
}
```

16. UpdateTask
>作用：更新指定的环境，负责更新环境的信息（环境名称；更新时间；更新版本等等）  
作用范围：服务端  
使用方法：
```
{
    "__type__": "UpdateTask"
}
```

### 条件任务
1. IfElse
>作用：做一个逻辑判断，可以配置条件，满足条件的任务列表；不满足条件的任务列表  
作用范围：服务端；客户端  
使用方法：
```
{
    "__type__": "IfElse",
    "Result": 0,
    "Logic": "==",
    "Cond": { //条件
        "FilePath": "/data/home/payneliu/services/apache-2.4/conf/httpd.conf",
        "SubString": "ServerName ${branch}.qgame.qq.com",
        "__type__": "ContainsTask"
    },
    "Body": [ //满足条件执行的任务列表
        {
            "FilePath": "/data/home/payneliu/services/apache-2.4/conf/httpd.conf",
            "Encoding": "utf8",
            "Replacer": [{
                "Partten": "###VirtualHostPlaceholder###",
                "Repl": "###${branch}_begin###\n<VirtualHost *>\nSetEnv APP_ENV dev\nDocumentRoot ${branchPath}/public/\nServerName ${branch}.qgame.qq.com\nErrorLog logs/${branch}.qgame.qq.com-error_log\nCustomLog logs/${branch}.qgame.qq.com-access_log common\n<Directory ${branchPath}/public/>\nOptions FollowSymLinks \nAllowOverride All\n#Order allow,deny \n#Allow from all\n</Directory>\n</VirtualHost>\n###${branch}_end###\n\n###VirtualHostPlaceholder###"
            }],
            "__type__": "ReplaceTask"
        }
    ],
    "ElseTask": [] //不满足条件执行的任务列表
}
```

### 还缺少功能
1. 对于vender的处理