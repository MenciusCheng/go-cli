package cmd

import (
	"github.com/MenciusCheng/go-cli/util/strategy"
	"github.com/MenciusCheng/go-cli/util/strategy/code_strategy"
	"github.com/spf13/cobra"
	"strconv"
)

var codeArgs = struct {
	FileDir              string
	FilePath             string
	SelectionStartLine   string
	SelectionEndLine     string
	SelectionStartColumn string
	SelectionEndColumn   string
	SelectedText         string
	FileText             string
	DeepseekApiKey       string
	QwenApiKey           string
}{}

var codeCmd = &cobra.Command{
	Use:   "code [prompt]",
	Short: "使用给定参数执行代码补全策略",
	Long:  `使用提示词和可选参数（如文件目录、文件路径等）执行代码补全策略`,
	Args:  cobra.MaximumNArgs(1), // 最多接受一个位置参数（prompt）
	RunE: func(cmd *cobra.Command, args []string) error {
		var prompt string
		if len(args) > 0 {
			prompt = args[0]
		}
		return codeHandler(prompt)
	},
}

func init() {
	// 添加所有可选的 flag 参数
	codeCmd.Flags().StringVar(&codeArgs.FileDir, "fileDir", "", "文件目录路径")
	codeCmd.Flags().StringVar(&codeArgs.FilePath, "filePath", "", "文件路径")
	codeCmd.Flags().StringVar(&codeArgs.SelectionStartLine, "selectionStartLine", "", "选择开始行号")
	codeCmd.Flags().StringVar(&codeArgs.SelectionEndLine, "selectionEndLine", "", "选择结束行号")
	codeCmd.Flags().StringVar(&codeArgs.SelectionStartColumn, "selectionStartColumn", "", "选择开始列号")
	codeCmd.Flags().StringVar(&codeArgs.SelectionEndColumn, "selectionEndColumn", "", "选择结束列号")
	codeCmd.Flags().StringVar(&codeArgs.SelectedText, "selectedText", "", "选中的文本内容")
	codeCmd.Flags().StringVar(&codeArgs.FileText, "fileText", "", "完整文件文本内容")
	codeCmd.Flags().StringVar(&codeArgs.DeepseekApiKey, "deepseekApiKey", "", "deepseek api key")
	codeCmd.Flags().StringVar(&codeArgs.QwenApiKey, "qwenApiKey", "", "qwen api key")

	rootCmd.AddCommand(codeCmd)
}

func codeHandler(prompt string) error {
	e := &strategy.Event{
		Prompt:               prompt,
		FileDir:              codeArgs.FileDir,
		FilePath:             codeArgs.FilePath,
		SelectionStartLine:   parseIntOrDefault(codeArgs.SelectionStartLine, 0),
		SelectionEndLine:     parseIntOrDefault(codeArgs.SelectionEndLine, 0),
		SelectionStartColumn: parseIntOrDefault(codeArgs.SelectionStartColumn, 0),
		SelectionEndColumn:   parseIntOrDefault(codeArgs.SelectionEndColumn, 0),
		SelectedText:         codeArgs.SelectedText,
		FileText:             codeArgs.FileText,
		DeepseekApiKey:       codeArgs.DeepseekApiKey,
		QwenApiKey:           codeArgs.QwenApiKey,
	}

	sm := strategy.NewStrategyManager()
	sm.RegisterStrategies(
		code_strategy.NewCodeStrategy(),
		strategy.NewEchoStrategy(),
	)

	return sm.HandleEvent(e)
}

// parseIntOrDefault 解析字符串为整数，如果为空或解析失败则返回默认值
func parseIntOrDefault(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return defaultValue
}
