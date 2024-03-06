package gpt

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
)

var (
	assistantId  = "asst_Yl3RYaM1TLc7RKWBDOw16Q3p"
	assistant    *openai.Assistant
	defaultModel = openai.GPT4TurboPreview
)

func assistantSetup() {
	resp, err := client.RetrieveAssistant(context.Background(), assistantId)
	if err != nil {
		panic(err)
	}
	assistant = &resp
}

func UploadFile(ctx context.Context, filename string) (string, error) {
	file, err := client.CreateFile(ctx, openai.FileRequest{
		FileName: filename,
		FilePath: filename,
		Purpose:  "assistants",
	})
	if err != nil {
		return "", err
	}
	return file.ID, nil
}

func DeleteFile(ctx context.Context, fileID string) error {
	err := client.DeleteFile(ctx, fileID)
	return err
}

type Message struct {
	Words   string
	FileIDs []string
}

func CreateThreadAndRun(ctx context.Context, messages []*Message) (string, string, error) {
	if messages == nil || len(messages) == 0 {
		return "", "", errors.New("messages cannot be nil or empty")
	}

	threadMessages := make([]openai.ThreadMessage, 0, len(messages))
	for index, message := range messages {
		if index%2 == 1 {
			threadMessages = append(threadMessages, openai.ThreadMessage{
				Role:    "assistant",
				Content: message.Words,
				FileIDs: message.FileIDs,
			})
		} else {
			threadMessages = append(threadMessages, openai.ThreadMessage{
				Role:    "user",
				Content: message.Words,
				FileIDs: message.FileIDs,
			})
		}
	}

	resp, err := client.CreateThreadAndRun(ctx, openai.CreateThreadAndRunRequest{
		RunRequest: openai.RunRequest{
			AssistantID: assistantId,
			Model:       &defaultModel,
		},
		Thread: openai.ThreadRequest{
			Messages: threadMessages,
		},
	})
	if err != nil {
		return "", "", err
	}
	return resp.ThreadID, resp.ID, nil
}

func AddThreadMessageAndRun(ctx context.Context, threadID string, message *Message) (string, error) {
	_, err := client.CreateMessage(ctx, threadID, openai.MessageRequest{
		Role:    "user",
		Content: message.Words,
		FileIds: message.FileIDs,
	})
	if err != nil {
		return "", err
	}
	resp, err := client.CreateRun(ctx, threadID, openai.RunRequest{
		AssistantID: assistantId,
	})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

var (
	defaultMessageLimit  = 2
	defaultMessageOrder  = "desc"
	defaultMessageBefore = ""
	defaultMessageAfter  = ""
)

func GetRunResult(ctx context.Context, threadID string, runID string) (string, bool, error) {
	resp, err := client.RetrieveRun(ctx, threadID, runID)
	if err != nil {
		return "", false, err
	}
	if resp.Status == openai.RunStatusCompleted {
		messages, err := client.ListMessage(ctx, threadID, &defaultMessageLimit, &defaultMessageOrder, &defaultMessageAfter, &defaultMessageBefore)
		if err != nil {
			return "", true, err
		}
		for _, item := range messages.Messages {
			if item.Role != "user" {
				return item.Content[0].Text.Value, true, err
			}
		}
		return "", true, err
	}
	return "", false, err
}
