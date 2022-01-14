package config

import (
	"github.com/fatih/color"
	"github.com/flytam/filenamify"
	"gitlab-api-user-enum-exploit/pkg/file_util"
	"net/url"
	"path"
	"strings"
)

type Config struct {

	// -------------------------------------------- 输入相关 ------------------------------------------------------

	// 直接给出能够枚举用户信息的url，剩下的版本信息之类的
	// 比如：http://www.foobar.com/api/v4/users/1，会自动识别后面的1，替换为递增，尝试拖用户名
	// 比如：http://www.foobar.com/api/v4/users，会在后面拼接递增的id尝试拖用户名
	ApiUrl string

	// 域名，比如： http://www.foobar.com
	Site string

	// 用来做批量测试，每行一个api url或者host
	InputFilePath string

	// -------------------------------------------- 请求相关 ------------------------------------------------------

	// 当遇到几个连续的404时，才认为用户是脱完了
	// 默认是10个
	Cutoff int

	// 请求时使用的代理
	Proxy string

	// 请求失败重试次数
	RequestMaxTryTimes int

	// -------------------------------------------- 输出相关 ------------------------------------------------------

	// 结果保存到指定的jsonline文件中，一个用户一个行json
	OutputJsonLineFile string

	// 只把处在active状态的用户的用户名单独导入到一个文件，用于做后续爆破相关的事情
	OutputUsernameFile string

	// 不指定jsonline路径，写入跟域名相同名称的jsonline文件
	OutputJsonLineDomainAuto bool
}

var RunConfig = &Config{
}

func ProcessConfig() ([]*Config, error) {
	// 对url做统一处理
	configSlice := make([]*Config, 0)

	// 输入来源1： 从文件中读取
	if RunConfig.InputFilePath != "" {
		lines, err := file_util.ReadLines(RunConfig.InputFilePath)
		if err != nil {
			color.Red("读取文件%s时发生错误： %s", RunConfig.InputFilePath, err.Error())
		} else {
			for _, line := range lines {
				if strings.Contains(strings.ToLower(line), "/api/") {
					// 是 https://gitlab.hzleaper.com:81/api/v4/users/200 的形式
					// 是 https://gitlab.hzleaper.com:81/api/v4/users/ 的形式
					// 是 https://gitlab.hzleaper.com:81/api/v4/users 的形式
					apiUrl := ProcessApiUrl(line)
					configSlice = append(configSlice, &Config{
						ApiUrl: apiUrl,
					})
				} else {
					// 认为是： https://gitlab.hzleaper.com:81 的形式，版本默认为v4
					apiUrl := path.Join(line, "/api/v4/users/")
					configSlice = append(configSlice, &Config{
						ApiUrl: apiUrl,
					})
				}
			}
		}
	}

	// 输入来源2： 从命令行直接指定的域名
	if RunConfig.Site != "" {
		apiUrl := path.Join(RunConfig.Site, "/api/v4/users/")
		configSlice = append(configSlice, &Config{
			ApiUrl: apiUrl,
		})
	}

	// 输入来源3： 从命令行直接指定的api url
	if RunConfig.ApiUrl != "" {
		configSlice = append(configSlice, &Config{
			ApiUrl: ProcessApiUrl(RunConfig.ApiUrl),
		})
	}

	// 输出文件修正，以及其他字段填充
	for _, config := range configSlice {

		// 输出相关
		config.OutputUsernameFile = RunConfig.OutputUsernameFile
		config.OutputJsonLineFile = RunConfig.OutputJsonLineFile
		if RunConfig.OutputJsonLineDomainAuto {
			parse, err := url.Parse(config.ApiUrl)
			if err != nil {
				color.Red("非法的API URL: %s", config.ApiUrl)
				continue
			}
			output, err := filenamify.Filenamify(parse.Host, filenamify.Options{})
			if err != nil {
				color.Red("从域名生成文件名错误，域名：%s，错误信息：%s", parse.Host, err.Error())
				continue
			}
			config.OutputJsonLineFile = "./output/" + output + ".jsonl"
		}

		// 其他子段
		config.Proxy = RunConfig.Proxy
		config.Cutoff = RunConfig.Cutoff
		config.RequestMaxTryTimes = RunConfig.RequestMaxTryTimes
	}

	// 把每个config执行一遍即可...
	return configSlice, nil
}

// ProcessApiUrl 把API URL统一格式化为： https://gitlab.hzleaper.com:81/api/v4/users/ 的形式
// 输入的可能格式：
// https://gitlab.hzleaper.com:81/api/v4/users/200
// https://gitlab.hzleaper.com:81/api/v4/users/
// https://gitlab.hzleaper.com:81/api/v4/users
func ProcessApiUrl(apiUrl string) string {
	parse, _ := url.Parse(apiUrl)
	s := strings.TrimRight(parse.Path, "/")
	split := strings.Split(s, "/")
	if len(split) == 0 {
		return apiUrl
	}
	if strings.ToLower(split[len(split)-1]) != "users" {
		split = split[0 : len(split)-1]
	}
	return parse.Scheme + "://" + parse.Host + "/" + strings.Join(split, "/") + "/"
}
