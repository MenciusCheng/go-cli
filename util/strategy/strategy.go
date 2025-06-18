package strategy

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Event 事件
type Event struct {
	Prompt               string `json:"prompt"`               // 用户输入的提示信息
	FileDir              string `json:"fileDir"`              // 文件所在目录
	FilePath             string `json:"filePath"`             // 文件完整路径
	SelectionStartLine   int    `json:"selectionStartLine"`   // 选中文本开始行号
	SelectionEndLine     int    `json:"selectionEndLine"`     // 选中文本结束行号
	SelectionStartColumn int    `json:"selectionStartColumn"` // 选中文本开始列号
	SelectionEndColumn   int    `json:"selectionEndColumn"`   // 选中文本结束列号
	SelectedText         string `json:"selectedText"`         // 选中的文本内容
	FileText             string `json:"fileText"`             // 完整文件内容
	DeepseekApiKey       string `json:"deepseekApiKey"`       // deepseek api key
	QwenApiKey           string `json:"qwenApiKey"`           // qwen api key
	RefStruct            string `json:"refStruct"`            // 参考结构体定义
}

func (e *Event) ToMapByJSON() map[string]interface{} {
	jsonData, _ := json.Marshal(e)
	result := make(map[string]interface{})
	_ = json.Unmarshal(jsonData, &result)
	return result
}

// Strategy 策略接口
type Strategy interface {
	// CanHandle 判断是否能处理该事件
	CanHandle(e *Event) bool
	// Handle 处理事件
	Handle(e *Event) error
	// GetName 获取策略名称
	GetName() string
}

// StrategyManager 策略管理器
type StrategyManager struct {
	strategies []Strategy
}

// NewStrategyManager 创建策略管理器
func NewStrategyManager() *StrategyManager {
	return &StrategyManager{
		strategies: make([]Strategy, 0),
	}
}

// RegisterStrategy 注册策略
func (sm *StrategyManager) RegisterStrategy(strategy Strategy) {
	sm.strategies = append(sm.strategies, strategy)
}

// RegisterStrategies 批量注册策略
func (sm *StrategyManager) RegisterStrategies(strategies ...Strategy) {
	for _, strategy := range strategies {
		sm.RegisterStrategy(strategy)
	}
}

// HandleEvent 处理事件
func (sm *StrategyManager) HandleEvent(event *Event) error {
	// 预处理事件
	if err := sm.preprocess(event); err != nil {
		return fmt.Errorf("preprocess failed: %w", err)
	}

	// 遍历所有策略，找到第一个能处理该事件的策略并执行
	for _, strategy := range sm.strategies {
		if strategy.CanHandle(event) {
			return strategy.Handle(event)
		}
	}

	// 如果没有找到合适的策略，返回错误
	return fmt.Errorf("no strategy found to handle the event")
}

// preprocess 预处理事件，补充文件内容和选中文本
func (sm *StrategyManager) preprocess(event *Event) error {

	if event.DeepseekApiKey == "" {
		// 检查API环境变量
		apiKey := os.Getenv("DEEPSEEK_API_KEY")
		if apiKey != "" {
			event.DeepseekApiKey = apiKey
		}
	}

	// 检查是否需要预处理
	if event.FilePath == "" || event.SelectionStartLine <= 0 || event.SelectionEndLine <= 0 {
		return nil
	}

	// 如果 FileText 或 SelectedText 不为空，则不需要预处理
	if event.FileText != "" || event.SelectedText != "" {
		return nil
	}

	// 读取文件内容
	fileContent, err := os.ReadFile(event.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", event.FilePath, err)
	}

	// 将文件内容转换为字符串
	fileText := string(fileContent)

	// 按行分割，保持原始的换行符和空白行
	lines := strings.Split(fileText, "\n")

	// 如果 FileText 为空，设置完整文件内容
	if event.FileText == "" {
		event.FileText = fileText
	}

	// 如果 SelectedText 为空，提取选中的行内容
	if event.SelectedText == "" &&
		(event.SelectionStartLine < event.SelectionEndLine || event.SelectionStartLine == event.SelectionEndLine && event.SelectionStartColumn < event.SelectionEndColumn) {
		// 验证行号范围
		if event.SelectionStartLine > len(lines) || event.SelectionEndLine > len(lines) {
			return fmt.Errorf("selection line numbers out of range: file has %d lines, but selection is from line %d to %d",
				len(lines), event.SelectionStartLine, event.SelectionEndLine)
		}

		if event.SelectionStartLine > event.SelectionEndLine {
			return fmt.Errorf("invalid selection: start line %d is greater than end line %d",
				event.SelectionStartLine, event.SelectionEndLine)
		}

		// 提取选中的行（注意：SelectionStartLine 从1开始，所以需要减1）
		selectedLines := lines[event.SelectionStartLine-1 : event.SelectionEndLine]
		event.SelectedText = strings.Join(selectedLines, "\n")
	}

	return nil
}
