package ask_strategy

import (
	"fmt"
	"github.com/MenciusCheng/go-cli/util/openai"
	"github.com/MenciusCheng/go-cli/util/strategy"
	"os"
)

func NewAskAnyStrategy() strategy.Strategy {
	return &AskAnyStrategy{}
}

type AskAnyStrategy struct {
}

func (s *AskAnyStrategy) CanHandle(e *strategy.Event) bool {
	return true
}

func (s *AskAnyStrategy) Handle(e *strategy.Event) error {
	if e.DeepseekApiKey == "" {
		return fmt.Errorf("apiKey为空")
	}

	// 检查API密钥
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("未设置 DEEPSEEK_API_KEY 环境变量")
	}
	fmt.Printf("正在咨询大模型...\n\n")
	client := openai.NewClient(e.DeepseekApiKey)
	err := client.StreamCodeAskWithPrompt(e.Prompt, func(token string) {
		// 流式打印内容
		fmt.Print(token)
	})
	if err != nil {
		return fmt.Errorf("咨询大模型失败: %v", err)
	}
	fmt.Println() // 在回答结束后添加换行
	return nil
}

func (s *AskAnyStrategy) GetName() string {
	return "AskAnyStrategy"
}
