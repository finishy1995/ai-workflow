package core

import (
	"fmt"
	"github.com/finishy1995/effibot-core/ai/pb/ai"
	"github.com/finishy1995/go-library/storage"
	"strings"
)

// Task AI Agent 任务
type Task struct {
	storage.Model
	ErrorCode          int32
	Role               Role
	Status             uint8
	TaskId             string `dynamo:",hash"`
	UserId             string
	WorkflowInstanceId string
	Input              []string
	Output             []string
	Param              map[string]interface{}
}

const (
	AIKeyContent            = "content"            // 	文字内容（脑图模式最终生成的文字，写/改点什么的文字内容）
	AIKeySelectedWords      = "selectedWords"      // 	用户选中的内容（写/改点什么）
	AIKeySelectedStartIndex = "selectedStartIndex" // 	用户选中的内容起始下标（写/改点什么）
	AIKeySelectedEndIndex   = "selectedEndIndex"   // 	用户选中的内容终止下标（写/改点什么）
	AIKeyInfo               = "info"               // 	写点什么用户填写的描述
	AIKeyTips               = "tips"               // 	修改内容时的用户需求（以 markdown 1. *\n2. *的格式传递）
	AIKeyBrain              = "brain"              // 	脑图内容
	AIKeyTitle              = "title"              // 	文档标题
	AIKeyOldContent         = "oldContent"         // 	在做出当前 AI 请求前，文字内容是什么
	AIKeyOldBrain           = "oldBrain"           // 	在做出当前 AI 请求前，脑图内容是什么
	AIKeyQuestion           = "question"           // 	在脑图文档中，使用 AI 问答功能用户提出的问题是什么
	AIKeyStyle              = "style"              //	风格
	AIKeyWordsNumber        = "word_count"         // 	字数
	AIKeyTone               = "tone"               //	语气
)

func KeyValue2Map(params []*ai.KeyValue, settings []string) map[string]interface{} {
	if params == nil {
		return nil
	}

	m := make(map[string]interface{}, len(params))
	for _, p := range params {
		if p.Key == AIKeyBrain {
			m[p.Key] = strings.Replace(p.Value, "\r\n", "\n", -1)
			continue
		}

		m[p.Key] = p.Value
	}
	if settings != nil && len(settings) > 0 {
		m[AIKeyStyle] = settings[0]
	}
	return m
}

func Map2KeyValue(params map[string]interface{}) []*ai.KeyValue {
	if params == nil {
		return nil
	}

	kvs := make([]*ai.KeyValue, 0, len(params))
	for k, v := range params {
		kvs = append(kvs, &ai.KeyValue{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		})
	}
	return kvs
}
