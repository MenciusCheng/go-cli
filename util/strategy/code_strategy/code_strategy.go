package code_strategy

import (
	"fmt"
	"github.com/MenciusCheng/go-cli/rule"
	"github.com/MenciusCheng/go-cli/util/openai"
	"github.com/MenciusCheng/go-cli/util/renderer"
	"github.com/MenciusCheng/go-cli/util/strategy"
	"os"
	"strings"
)

func NewCodeStrategy() strategy.Strategy {
	return &CodeStrategy{}
}

type CodeStrategy struct {
}

func (s *CodeStrategy) CanHandle(e *strategy.Event) bool {
	// 必须选中代码
	if e.SelectedText == "" {
		return false
	}
	return true
}

func (s *CodeStrategy) Handle(e *strategy.Event) error {
	fmt.Printf("策略名称: %s\n", s.GetName())
	fmt.Printf("任意代码补全\n")

	if e.DeepseekApiKey == "" {
		return fmt.Errorf("apiKey为空")
	}

	eMap := e.ToMapByJSON()
	if len([]rune(e.FileText)) > 100000 {
		// 截取字符串
		eMap["fileText"] = string([]rune(e.FileText)[:100000])
	} else {
		eMap["fileText"] = e.FileText
	}

	render := renderer.New()
	tmplContent := rule.CodeRuleTemplate
	// 渲染模板
	prompt, err := render.RenderString(tmplContent, eMap)
	if err != nil {
		return fmt.Errorf("渲染模板失败: %w", err)
	}

	fmt.Println("\n=== 提示词 ===")
	fmt.Printf("%s\n", prompt)

	fmt.Println("\n=== 正在执行代码补全 ===")
	client := openai.NewClient(e.DeepseekApiKey)
	var completedCode strings.Builder
	err = client.StreamCodeCompletionWithPrompt(prompt, func(token string) {
		// 流式打印内容
		fmt.Print(token)
		// 同时将内容写入到 builder 中
		completedCode.WriteString(token)
	})
	if err != nil {
		return err
	}
	completedCodeStr := client.TrimMarkdown(completedCode.String())

	fmt.Println("\n=== 正在替换代码 ===")
	err = replaceCodeInFile(e.FilePath, e.FileText, e.SelectionStartLine, e.SelectionEndLine, completedCodeStr)
	if err != nil {
		return fmt.Errorf("替换代码失败: %v", err)
	}

	fmt.Printf("代码补全完成，已更新文件: %s\n", e.FilePath)
	return nil
}

func (s *CodeStrategy) GetName() string {
	return "CodeStrategy"
}

func replaceCodeInFile(filePath string, fileText string, startLine, endLine int, completedCode string) error {
	allLines := strings.Split(fileText, "\n")

	var newLines []string
	completedLines := strings.Split(completedCode, "\n")

	if startLine == 0 && endLine == 0 {
		// 替换整个文件
		newLines = completedLines
	} else {
		// 替换指定范围的行
		// 添加开始行之前的内容
		if startLine > 1 {
			newLines = append(newLines, allLines[:startLine-1]...)
		}

		// 添加补全后的代码
		newLines = append(newLines, completedLines...)

		// 添加结束行之后的内容
		if endLine < len(allLines) {
			newLines = append(newLines, allLines[endLine:]...)
		}
	}

	// 写入新内容
	newContent := strings.Join(newLines, "\n")
	err := os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}
