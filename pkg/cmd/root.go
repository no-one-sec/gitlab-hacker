package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "gitlab-api-user-enum-exploit",
	Short: "gitlab-api-user-enum-exploit",
	Long:  `gitlab api user enum exploit tool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gitlab api user enum exploit tool")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
