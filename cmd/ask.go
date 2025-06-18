package cmd

import (
	"github.com/MenciusCheng/go-cli/util/strategy"
	"github.com/MenciusCheng/go-cli/util/strategy/ask_strategy"
	"github.com/spf13/cobra"
	"strings"
)

var askArgs = struct {
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

var askCmd = &cobra.Command{
	Use:   "ask [prompt]",
	Short: "咨询大模型问题",
	Long:  `向大模型提问并获得回答。必须提供问题作为参数。`,
	Args:  cobra.MinimumNArgs(1), // 至少需要一个参数
	RunE: func(cmd *cobra.Command, args []string) error {
		// 将所有参数连接成一个问题
		prompt := strings.Join(args, " ")
		return askHandler(prompt)
	},
}

func init() {
	// 添加所有可选的 flag 参数
	askCmd.Flags().StringVar(&askArgs.FileDir, "fileDir", "", "文件目录路径")
	askCmd.Flags().StringVar(&askArgs.FilePath, "filePath", "", "文件路径")
	askCmd.Flags().StringVar(&askArgs.SelectionStartLine, "selectionStartLine", "", "选择开始行号")
	askCmd.Flags().StringVar(&askArgs.SelectionEndLine, "selectionEndLine", "", "选择结束行号")
	askCmd.Flags().StringVar(&askArgs.SelectionStartColumn, "selectionStartColumn", "", "选择开始列号")
	askCmd.Flags().StringVar(&askArgs.SelectionEndColumn, "selectionEndColumn", "", "选择结束列号")
	askCmd.Flags().StringVar(&askArgs.SelectedText, "selectedText", "", "选中的文本内容")
	askCmd.Flags().StringVar(&askArgs.FileText, "fileText", "", "完整文件文本内容")
	askCmd.Flags().StringVar(&askArgs.DeepseekApiKey, "deepseekApiKey", "", "deepseek api key")
	askCmd.Flags().StringVar(&askArgs.QwenApiKey, "qwenApiKey", "", "qwen api key")

	rootCmd.AddCommand(askCmd)
}

func askHandler(prompt string) error {
	e := &strategy.Event{
		Prompt:               prompt,
		FileDir:              askArgs.FileDir,
		FilePath:             askArgs.FilePath,
		SelectionStartLine:   parseIntOrDefault(askArgs.SelectionStartLine, 0),
		SelectionEndLine:     parseIntOrDefault(askArgs.SelectionEndLine, 0),
		SelectionStartColumn: parseIntOrDefault(askArgs.SelectionStartColumn, 0),
		SelectionEndColumn:   parseIntOrDefault(askArgs.SelectionEndColumn, 0),
		SelectedText:         askArgs.SelectedText,
		FileText:             askArgs.FileText,
		DeepseekApiKey:       askArgs.DeepseekApiKey,
		QwenApiKey:           askArgs.QwenApiKey,
	}

	sm := strategy.NewStrategyManager()
	sm.RegisterStrategies(
		ask_strategy.NewAskCodeStrategy(),
		ask_strategy.NewAskAnyStrategy(),
		strategy.NewEchoStrategy(),
	)

	return sm.HandleEvent(e)
}
