package tools

import (
  openai "github.com/sashabaranov/go-openai"
  "encoding/json"
  "strings"
  "strconv"
  "fmt"
)

type ToolType string
//type InputArgs struct {
//    Numbers []int `json:"numbers"`
//}

type InputArgs struct {
    // 使用 json.RawMessage 延迟解析，或自定义类型
    Numbers []int `json:"-"` // 忽略默认解析
}

// UnmarshalJSON 自定义反序列化逻辑
func (a *InputArgs) UnmarshalJSON(data []byte) error {
    // 先尝试标准解析
    var temp struct {
        Numbers json.RawMessage `json:"numbers"` // 原始 JSON 数据
    }
    if err := json.Unmarshal(data, &temp); err != nil {
        return err
    }

    // 尝试解析为 []int（正确格式）
    var intArray []int
    if err := json.Unmarshal(temp.Numbers, &intArray); err == nil {
        a.Numbers = intArray
        return nil
    }

    // 尝试解析为字符串（错误格式 "[1, 2]"）
    var str string
    if err := json.Unmarshal(temp.Numbers, &str); err == nil {
        // 去除可能的引号和括号，解析数字
        str = strings.Trim(str, `"[] `)
        parts := strings.Split(str, ",")
        for _, p := range parts {
            p = strings.TrimSpace(p)
            if p == "" {
                continue
            }
            num, err := strconv.Atoi(p)
            if err != nil {
                return fmt.Errorf("无法解析数字: %s", p)
            }
            a.Numbers = append(a.Numbers, num)
        }
        return nil
    }

    return fmt.Errorf("无法解析 numbers 字段: %s", string(temp.Numbers))
}

const (
    ToolTypeFunction ToolType = "function"
)


type Tool struct {
    Type     ToolType            `json:"type"`
    Function *FunctionDefinition `json:"function,omitempty"`
}


type FunctionDefinition struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    Parameters any `json:"parameters"`
}

var AddToolDefine = openai.Tool{
    Type: "function",
    Function: &openai.FunctionDefinition{
        Name: "AddTool",
        Description: `
        Use this tool for addition calculations.
            example:
                1+2 =?
            then Action Input is: 1,2
        `,
        Parameters: `{"type":"object","properties":{"numbers":{"type":"array","items":{"type":"integer"}}}}`,
    },
}

func AddTool(numbers []int) int {
    sum := 0
    for _, n := range numbers {
        sum += n
    }
    return sum
}
