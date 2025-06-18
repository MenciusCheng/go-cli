package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// 版本信息常量
const (
	Version = "0.1.0"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return versionHandler()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionHandler() error {
	fmt.Printf("版本: %s\n", Version)
	return nil
}
