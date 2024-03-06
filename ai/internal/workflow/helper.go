package workflow

import (
	"fmt"
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/ai/pb/ai"
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
	"github.com/finishy1995/effibot-core/library/id"
	"github.com/finishy1995/go-library/routine"
	"github.com/finishy1995/go-library/storage"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

var (
	type2Workflow = map[ai.NodeType]string{
		ai.NodeType_ContentCreate:          "wf_content_create",
		ai.NodeType_ContentTips:            "wf_content_tips",
		ai.NodeType_BrainCreate:            "wf_brain_create",
		ai.NodeType_BrainTips:              "wf_brain_tips",
		ai.NodeType_ContentUpdate:          "wf_content_update",
		ai.NodeType_ContentModify:          "wf_content_modify",
		ai.NodeType_ContentUpdateByBrain:   "wf_content_update_by_brain",
		ai.NodeType_ContentTranslate:       "wf_content_translate",
		ai.NodeType_WordsDescribe:          "wf_words_describe",
		ai.NodeType_BrainUpdate:            "wf_brain_update",
		ai.NodeType_BrainModify:            "wf_brain_modify",
		ai.NodeType_BrainUpdateByContent:   "wf_brain_update_by_content",
		ai.NodeType_BrainExpand:            "wf_brain_expand",
		ai.NodeType_BrainExpandToNewBrain:  "wf_brain_expand_to_new_brain",
		ai.NodeType_BrainQuestionAsk:       "wf_brain_question_ask",
		ai.NodeType_ContentCreateByBrain:   "wf_content_create_by_brain",
		ai.NodeType_ImageCreateByBrain:     "wf_image_create_by_brain",
		ai.NodeType_BrainCreateByContent:   "wf_brain_create_by_content",
		ai.NodeType_TempContentPersonalize: "wf_temp_content_personalize",
	}
	workflow2Type = map[string]ai.NodeType{}

	dbInst storage.Storage
)

func init() {
	for key, value := range type2Workflow {
		workflow2Type[value] = key
	}
}

func Init(st storage.Storage) {
	dbInst = st
	err := st.CreateTable(Config{}, base.WorkflowTableName)
	if err != nil {
		panic(err)
	}
	err = data.SetupPool(base.WorkflowDataName, NewBaseAiWorkflow, &Config{}, &data.LoaderConfig{
		Type:      data.StorageWithFileInit,
		FilePath:  base.DataPath + DataPath,
		TableName: base.WorkflowTableName,
	})
	if err != nil {
		panic(err)
	}
}

func Process(userId, workflowId string, param map[string]interface{}) (taskId string, err error) {
	item, ok := data.GetOneItem(base.WorkflowDataName, workflowId)
	if !ok {
		return "", ErrInvalidWorkflowId
	}
	workflow, ok := item.(AiWorkflow)
	if !ok {
		return "", ErrInvalidWorkflowId
	}
	task := &core.Task{
		ErrorCode:          base.CodeOk,
		TaskId:             id.GenerateID().String(),
		UserId:             userId,
		WorkflowInstanceId: workflowId,
		Input:              []string{},
		Output:             nil,
		Param:              param,
		Status:             uint8(ai.NodeStatus_Running),
	}

	if err = dbInst.Create(*task, ""); err != nil {
		return "", err
	}

	err = routine.Run(true, func() {
		workflow.Process(task)

		// read
		var readTask core.Task
		if errRoutine := dbInst.First(&readTask, "", task.TaskId); errRoutine != nil {
			logx.Errorf("errRoutine task get error, task id: %s, ai result: %+v, error: %s",
				task.TaskId, task.Output, errRoutine.Error())
			return
		}

		// save
		readTask.ErrorCode = task.ErrorCode
		readTask.Input = task.Input
		readTask.Output = task.Output
		readTask.Param = task.Param
		readTask.Status = uint8(ai.NodeStatus_Ready)
		if errRoutine := dbInst.Save(&readTask, ""); errRoutine != nil {
			logx.Errorf("errRoutine task save error, task id: %s, ai result: %+v, error: %s",
				task.TaskId, task.Output, errRoutine.Error())
		}
	})

	return task.TaskId, err
}

// ProcessOld Deprecated
func ProcessOld(typ ai.NodeType, oldSessionId, userId string, param map[string]interface{}) (sessionId string, nodeId string, err error) {
	if oldSessionId == "" {
		sessionId = fmt.Sprintf("%s%s%s", userId, SessionIdMark, id.GenerateID().String())
	} else {
		sessionId = oldSessionId
		// 从 sessionId 中，获取 userId
		parts := strings.Split(sessionId, SessionIdMark)
		if len(parts) > 0 {
			userId = parts[0]
		} else {
			// Handle the error if sessionId format is not as expected
			err = ErrInvalidSessionId
			return
		}
	}
	nodeId, err = Process(userId, GetWorkflowId(typ), param)
	return
}

func GetWorkflowId(typ ai.NodeType) string {
	return type2Workflow[typ]
}

func GetType(workflowId string) ai.NodeType {
	return workflow2Type[workflowId]
}
