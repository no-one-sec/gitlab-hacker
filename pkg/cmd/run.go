package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gitlab-api-user-enum-exploit/pkg/config"
	"gitlab-api-user-enum-exploit/pkg/core"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Long:  `开始拉取用户信息`,
	Run: func(cmd *cobra.Command, args []string) {
		configSlice, err := config.ProcessConfig()
		if err != nil {
			color.Red("解析输入配置错误： %s", err.Error())
			return
		}
		for _, config := range configSlice {
			x := core.NewGitlabUserEnum(config)
			err := x.Init()
			if err != nil {
				color.Red("初始化环境错误：%s", err.Error())
				continue
			}
			x.Run()
		}
		color.BlackString("ALL DONE")
	},
	TraverseChildren: true,
}

func init() {

	// 输入相关
	runCmd.Flags().StringVar(&config.RunConfig.ApiUrl, "api-url", "", "有用户泄露的API地址，比如 https://foo.com/api/v4/users/1")
	runCmd.Flags().StringVar(&config.RunConfig.Site, "site", "", "所关联的网站，比如 https://foo.com")
	runCmd.Flags().StringVar(&config.RunConfig.InputFilePath, "from-file", "", "从文件中批量运行，文件中每行是一个api url或者domain")

	// 请求相关
	runCmd.Flags().IntVar(&config.RunConfig.Cutoff, "cut-off", 10, "遇到多少个连续的404时认为是拖完所有用户了")
	runCmd.Flags().StringVar(&config.RunConfig.Proxy, "proxy", "", "请求时使用代理")
	runCmd.Flags().IntVar(&config.RunConfig.RequestMaxTryTimes, "request-max-try-times", 3, "请求重试次数")

	// 输出相关
	runCmd.Flags().StringVar(&config.RunConfig.OutputJsonLineFile, "output-json-line-file", "", "把拉取到的所有用户的信息保存到一个JsonLine文件")
	runCmd.Flags().StringVar(&config.RunConfig.OutputUsernameFile, "output-username-file", "", "把active状态的用户单独保存到一个文件中")
	runCmd.Flags().BoolVar(&config.RunConfig.OutputJsonLineDomainAuto, "output-by-domain", false, "根据域名自动生成文件名保存")

	rootCmd.AddCommand(runCmd)
}
