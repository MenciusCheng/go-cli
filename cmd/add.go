package cmd

import (
	"fmt"
	"github.com/MenciusCheng/go-cli/templates"
	"github.com/MenciusCheng/go-cli/util/renderer"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "添加新命令到项目中",
	Long:  `添加新命令到项目中，将根据模板生成命令文件。命令名称必须指定且不能与现有命令重名。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return addHandler(args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addHandler(name string) error {
	// 检查命令名是否为空
	if name == "" {
		return fmt.Errorf("命令名称不能为空")
	}

	// 生成文件路径
	var filePath string
	cmdDir := "cmd"

	// 检查当前目录下是否存在 cmd 目录
	if info, err := os.Stat(cmdDir); err == nil && info.IsDir() {
		// cmd 目录存在，将文件生成到 cmd 目录下
		filePath = filepath.Join(cmdDir, fmt.Sprintf("%s.go", name))
	} else {
		// cmd 目录不存在，生成到当前目录
		filePath = fmt.Sprintf("%s.go", name)
	}

	// 检查文件是否已存在
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("命令文件 '%s' 已存在", filePath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("检查文件时出错: %v", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"name": name,
	}

	// 渲染模板并写入文件
	err := renderer.New().RenderToFile(templates.AddTemplate, data, filePath)
	if err != nil {
		return fmt.Errorf("创建命令文件失败: %v", err)
	}

	// 成功提示
	fmt.Printf("命令 '%s' 创建成功，文件路径: %s\n", name, filePath)

	return nil
}
