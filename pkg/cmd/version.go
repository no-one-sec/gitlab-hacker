package cmd

// 展示版本信息

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "v",
	Long:  `查看版本信息`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("当前版本为v0.1（2022-1-15 03:14:58）")
	},
	TraverseChildren: true,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
