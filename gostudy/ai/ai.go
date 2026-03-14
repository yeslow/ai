package ai

import (
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

var modelString string = "qwen3.5-plus" 
func NewOpenAiClient() *openai.Client {
	token := os.Getenv("ALICP_API_KEY")
  dashscope_url := "https://coding.dashscope.aliyuncs.com/v1"
  config := openai.DefaultConfig(token)
  config.BaseURL = dashscope_url


  return openai.NewClientWithConfig(config)
}

func Chat(message []openai.ChatCompletionMessage) openai.ChatCompletionMessage {
	c := NewOpenAiClient()
	rsp, err := c.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model:    "qwen3.5-plus",
		Messages: message,
		Stream:   false,
	})
	if err != nil {
		log.Println(err)
		return openai.ChatCompletionMessage{}
	}

	return rsp.Choices[0].Message
}

var MessageStore ChatMessages
type ChatMessages []openai.ChatCompletionMessage

func (cm *ChatMessages) AddFor(role string, msg string, toolCalls []openai.ToolCall) {
	*cm = append(*cm, openai.ChatCompletionMessage{
		Role:      role,
		Content:   msg,
		ToolCalls: toolCalls,
	})
}

func (cm *ChatMessages) ToMessage() []openai.ChatCompletionMessage {
	ret := make([]openai.ChatCompletionMessage, len(*cm))
	for index, c := range *cm {
		ret[index] = c
	}
	return ret
}

func (cm *ChatMessages) AddForTool(content, toolCallID string) {
	*cm = append(*cm, openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    content,
		ToolCallID: toolCallID,
	})
}

func ToolsChat(message []openai.ChatCompletionMessage, tools []openai.Tool) openai.ChatCompletionMessage {
    c := NewOpenAiClient()
    rsp, err := c.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
        Model:      modelString,
        Messages:   message,
        Tools:      tools,
        ToolChoice: "auto",
    })
    if err != nil {
        log.Println(err)
        return openai.ChatCompletionMessage{}
    }


    return rsp.Choices[0].Message
}
