package main

import (
	"fmt"
  "log"
  "agentstudy/tools"
	"agentstudy/ai"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
//	ai.MessageStore.AddFor(openai.ChatMessageRoleSystem, "你是一个足球领域的专家，请尽可能地帮我回答与足球相关的问题。")
//	ai.MessageStore.AddFor(openai.ChatMessageRoleUser, "C 罗是哪个国家的足球运动员？")
//	ai.MessageStore.AddFor(openai.ChatMessageRoleAssistant, "C 罗是葡萄牙足球运动员。")
//	ai.MessageStore.AddFor(openai.ChatMessageRoleUser, "内马尔呢？")

//	response := ai.Chat(ai.MessageStore.ToMessage())
//	fmt.Println(response.Content)

  toolsList := make([]openai.Tool, 0)
  toolsList = append(toolsList, tools.AddToolDefine)


  prompt := "1+2+3+4+5+6=?"
  ai.MessageStore.AddFor(openai.ChatMessageRoleUser, prompt, nil)


  response := ai.ToolsChat(ai.MessageStore.ToMessage(), toolsList)
  toolCalls := response.ToolCalls


  fmt.Println("大模型的回复是: ", response.Content)
  fmt.Println("大模型选择的工具是: ", toolCalls)

  if toolCalls != nil {
    var result int
    var args tools.InputArgs
    // err := json.Unmarshal([]byte(toolCalls[0].Function.Arguments), &args)
    err := args.UnmarshalJSON([]byte(toolCalls[0].Function.Arguments))
    if err != nil {
        log.Fatalln("json unmarshal err: ", err.Error())
    }


    if toolCalls[0].Function.Name == tools.AddToolDefine.Function.Name {
        result = tools.AddTool(args.Numbers)
    }


    fmt.Println("函数计算结果: ", result)
    ai.MessageStore.AddFor(openai.ChatMessageRoleAssistant, response.Content, response.ToolCalls)
    ai.MessageStore.AddForTool(fmt.Sprintf("%d", result), toolCalls[0].ID)


    response := ai.ToolsChat(ai.MessageStore.ToMessage(), toolsList)


    fmt.Println("大模型的最终回复: ", response.Content)
  }
}
