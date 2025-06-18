package strategy

import (
	"encoding/json"
	"fmt"
)

func NewEchoStrategy() Strategy {
	return &EchoStrategy{}
}

// 示例策略实现
type EchoStrategy struct{}

func (s *EchoStrategy) CanHandle(e *Event) bool {
	return true
}

func (s *EchoStrategy) Handle(e *Event) error {
	fmt.Printf("策略名称: %s\n", s.GetName())
	fmt.Printf("命中默认策略，事件参数：")
	// 将事件转换为 JSON 格式并打印
	eventJSON, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal event to JSON: %w", err)
	}
	fmt.Printf("%s\n", string(eventJSON))
	return nil
}

func (s *EchoStrategy) GetName() string {
	return "EchoStrategy"
}
