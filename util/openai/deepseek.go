package openai

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"strings"
)

type Client struct {
	client *openai.Client
	Model  string
}

func NewClient(authToken string) *Client {
	config := openai.DefaultConfig(authToken)
	config.BaseURL = "https://api.deepseek.com/v1"
	client := openai.NewClientWithConfig(config)
	return &Client{
		client: client,
		Model:  "deepseek-coder",
	}
}

// StreamCodeCompletionWithPrompt 流式代码补全，实时输出补全过程
func (c *Client) StreamCodeCompletionWithPrompt(prompt string, callback func(string)) error {
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model: c.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一个专业的代码补全助手。只返回代码，不要添加任何解释或markdown格式。",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.1,
		Stream:      true,
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("创建流式代码补全失败: %v", err)
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("接收流式代码补全数据失败: %v", err)
		}

		if len(response.Choices) > 0 {
			delta := response.Choices[0].Delta
			if delta.Content != "" {
				callback(delta.Content)
			}
		}
	}
	return nil
}

// 简单清理，移除可能的markdown代码块标记
func (c *Client) TrimMarkdown(code string) string {
	completedCode := strings.TrimSpace(code)
	if strings.HasPrefix(completedCode, "```") {
		lines := strings.Split(completedCode, "\n")
		if len(lines) > 2 {
			// 移除第一行的```language和最后一行的```
			completedCode = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}
	return completedCode
}

// StreamCodeAskWithPrompt 流式代码咨询
func (c *Client) StreamCodeAskWithPrompt(prompt string, callback func(string)) error {
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model: c.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一个专业的开发者，回答代码相关问题",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.1,
		Stream:      true,
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("创建流式代码补全失败: %v", err)
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("接收流式数据失败: %v", err)
		}

		if len(response.Choices) > 0 {
			delta := response.Choices[0].Delta
			if delta.Content != "" {
				callback(delta.Content)
			}
		}
	}
	return nil
}
