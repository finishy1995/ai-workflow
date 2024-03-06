package tools

import (
	"fmt"
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/sashabaranov/go-openai"
	lua "github.com/yuin/gopher-lua"
	"github.com/zeromicro/go-zero/core/logx"
)

func SetLuaData(l *lua.LState, task *core.Task, param map[string]interface{}) {
	// 创建一个新的 Lua 表
	paramsTable := l.NewTable()
	for key, value := range param {
		l.SetField(paramsTable, key, lua.LString(fmt.Sprintf("%v", value)))
	}
	// 将表设置为全局变量
	l.SetGlobal("AgentData", paramsTable)

	if task != nil {
		// 创建一个新的 Lua 表
		requestTable := l.NewTable()

		// 将 Request 对象的每个字段添加到 Lua 表中
		l.SetField(requestTable, "Role", lua.LNumber(task.Role))
		l.SetField(requestTable, "UserId", lua.LString(task.UserId))
		l.SetField(requestTable, "WorkflowInstanceId", lua.LString(task.WorkflowInstanceId))
		l.SetField(requestTable, "TaskId", lua.LString(task.TaskId))
		// 创建一个新的 Lua 表用于存储 Input 列表
		inputTable := l.NewTable()
		for _, input := range task.Input {
			inputTable.Append(lua.LString(input))
		}
		l.SetField(requestTable, "Input", inputTable)

		// 将 task.Params 添加到 Lua 表中
		paramsSubTable := l.NewTable()
		for key, value := range task.Param {
			l.SetField(paramsSubTable, key, lua.LString(fmt.Sprintf("%v", value)))
		}
		l.SetField(requestTable, "Param", paramsSubTable)

		// 将表设置为全局变量
		l.SetGlobal("AgentTask", requestTable)
	}
}

func SimpleGenerate(L *lua.LState) int {
	// 获取传递给 Lua 函数的参数数量
	top := L.GetTop()
	if top < 5 {
		L.Push(lua.LNil)                            // 第一个返回值为 nil
		L.Push(lua.LString("not enough arguments")) // 第二个返回值为错误信息
		return 2                                    // 返回两个值
	}

	system := L.ToString(1)
	promptTable := L.ToTable(2)
	// 初始化一个空的 Go 字符串切片来存储 prompt 数据
	var prompts []string
	// 遍历 Lua 表，并将每个元素添加到 prompts 切片中
	promptTable.ForEach(func(_ lua.LValue, value lua.LValue) {
		// 确保表中的元素是字符串
		if str, ok := value.(lua.LString); ok {
			prompts = append(prompts, string(str))
		}
	})
	model := openai.GPT3Dot5Turbo
	modelChoice := L.ToString(3)
	if modelChoice == "Advanced" {
		model = openai.GPT4TurboPreview
	}
	maxToken := L.ToInt(4)
	temperature := L.ToNumber(5)

	logx.Infof("system: %s, model: %s, maxToken: %d, temperature: %f", system, model, maxToken, temperature)
	// 检查 prompts 是否正确
	for _, prompt := range prompts {
		logx.Infof("prompt: %s", prompt)
	}

	//completion, err := gpt.ChatCompletion.CreateChatCompletion(
	//	context.Background(),
	//	system,
	//	prompts,
	//	gpt.WithModel(model),
	//	gpt.WithMaxToken(maxToken),
	//	gpt.WithTemperature(float32(temperature)),
	//)
	//if err != nil {
	//	logx.Errorf("gpt chat completion error: %s", err.Error())
	//	L.Push(lua.LNil)                 // 第一个返回值为 nil
	//	L.Push(lua.LString(err.Error())) // 第二个返回值为错误信息
	//	return 2                         // 返回两个值
	//}
	//// 成功时推送生成的文本和 nil
	//L.Push(lua.LString(completion)) // 第一个返回值为生成的文本
	//L.Push(lua.LNil)                // 第二个返回值为 nil

	L.Push(lua.LNil)                       // 第一个返回值为 nil
	L.Push(lua.LString("cancel gpt task")) // 第二个返回值为错误信息

	return 2 // 返回两个值
}
