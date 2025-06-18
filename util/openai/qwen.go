package openai

import (
	"github.com/sashabaranov/go-openai"
)

func NewQwenClient(authToken string) *Client {
	config := openai.DefaultConfig(authToken)
	config.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	client := openai.NewClientWithConfig(config)
	return &Client{
		client: client,
		Model:  "qwen-plus",
	}
}
