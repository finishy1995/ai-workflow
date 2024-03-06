package gpt

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type chatCompletion struct {
	defaultReq openai.ChatCompletionRequest
}

var ChatCompletion = &chatCompletion{
	defaultReq: openai.ChatCompletionRequest{
		Model:       openai.GPT4TurboPreview,
		MaxTokens:   4000,
		Temperature: 0.4,
		TopP:        1,
	},
}

func (c *chatCompletion) SetDefaultRequest(req openai.ChatCompletionRequest) {
	c.defaultReq = req
}

func (c *chatCompletion) CreateChatCompletion(ctx context.Context, system string, prompts []string, opts ...Option) (string, error) {
	req := c.defaultReq
	c.setOptions(&req, opts...)
	req.Messages = c.makeMessage(system, prompts)

	if usageInstance.checkRecent15min() {
		return "", fmt.Errorf("too many requests now, please try later")
	}

	start := time.Now()
	resp, err := client.CreateChatCompletion(ctx, req)
	logx.WithContext(ctx).Infof("openai request duration: %s", time.Since(start))
	if err != nil {
		return "", err
	}
	var results []string
	for _, choice := range resp.Choices {
		results = append(results, choice.Message.Content)
	}
	logx.WithContext(ctx).Infof("using token \tQuestion: %d Answer: %d Total: %d", resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	logx.WithContext(ctx).Debugf("receive OpenAI response:\t %+v", resp)

	usageInstance.update(resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)

	if len(results) == 0 {
		return "", fmt.Errorf("no results")
	}
	return GetMarkdown(results[0]), nil
}

func (c *chatCompletion) makeMessage(system string, prompts []string) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{{
		Role:    "system",
		Content: system,
	}}
	for i, prompt := range prompts {
		role := "user"
		if i%2 == 1 {
			role = "assistant"
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: prompt,
		})
	}
	return messages
}

func (c *chatCompletion) setOptions(req *openai.ChatCompletionRequest, opts ...Option) {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	if options.MaxToken != 0 {
		req.MaxTokens = options.MaxToken
	}

	if options.Model != "" {
		req.Model = options.Model
	}

	if options.Temperature != 0 {
		req.Temperature = options.Temperature
	}

	if options.TopP != 0 {
		req.TopP = options.TopP
	}

	if options.N != 0 {
		req.N = options.N
	}

	if options.PresencePenalty != 0 {
		req.PresencePenalty = options.PresencePenalty
	}

	if options.FrequencyPenalty != 0 {
		req.FrequencyPenalty = options.FrequencyPenalty
	}

	if options.User != "" {
		req.User = options.User
	}
}
