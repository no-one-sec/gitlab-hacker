package cmd

// 显示作者相关信息

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var authorCmd = &cobra.Command{
	Use:   "author",
	Short: "",
	Long:  `About Author`,
	Run: func(cmd *cobra.Command, args []string) {
		builder := strings.Builder{}
		builder.WriteString("\n\n")
		builder.WriteString("\t\tby CC11001100\n\n")
		builder.WriteString("\t\tGitHub: https://github.com/CC11001100\n")
		builder.WriteString("\t\tBlog: http://CC11001100.com\n\n")
		fmt.Println(builder.String())
	},
	TraverseChildren: true,
}

func init() {
	rootCmd.AddCommand(authorCmd)
}
