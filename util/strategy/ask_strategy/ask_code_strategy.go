package ask_strategy

import (
	"fmt"
	"github.com/MenciusCheng/go-cli/rule"
	"github.com/MenciusCheng/go-cli/util/openai"
	"github.com/MenciusCheng/go-cli/util/renderer"
	"github.com/MenciusCheng/go-cli/util/strategy"
)

func NewAskCodeStrategy() strategy.Strategy {
	return &AskCodeStrategy{}
}

type AskCodeStrategy struct {
}

func (s *AskCodeStrategy) CanHandle(e *strategy.Event) bool {
	// 必须选中代码
	if e.SelectedText == "" {
		return false
	}
	return true
}

func (s *AskCodeStrategy) Handle(e *strategy.Event) error {
	fmt.Printf("策略名称: %s\n", s.GetName())
	fmt.Printf("任意代码咨询\n")

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
	tmplContent := rule.AskRuleTemplate
	// 渲染模板
	prompt, err := render.RenderString(tmplContent, eMap)
	if err != nil {
		return fmt.Errorf("渲染模板失败: %w", err)
	}

	fmt.Println("\n=== 提示词 ===")
	fmt.Printf("%s\n", prompt)

	fmt.Println("\n=== 大模型回答 ===")
	client := openai.NewClient(e.DeepseekApiKey)
	err = client.StreamCodeAskWithPrompt(prompt, func(token string) {
		// 流式打印内容
		fmt.Print(token)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *AskCodeStrategy) GetName() string {
	return "AskCodeStrategy"
}
