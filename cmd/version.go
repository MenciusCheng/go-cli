package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// 版本信息变量，将在构建时通过 ldflags 注入
var (
	Version   = "dev"     // 版本号
	BuildTime = "unknown" // 构建时间
	GitCommit = "unknown" // Git 提交哈希
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示程序的详细版本信息，包括版本号、构建时间和 Git 提交哈希`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return versionHandler()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionHandler() error {
	fmt.Printf("版本: %s\n", Version)
	fmt.Printf("构建时间: %s\n", BuildTime)
	fmt.Printf("Git 提交: %s\n", GitCommit)
	return nil
}
