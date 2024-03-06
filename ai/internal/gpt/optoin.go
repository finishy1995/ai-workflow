package gpt

// Option 配置项
type Option func(*Options)

// Options 配置结构体
type Options struct {
	MaxToken         int
	Model            string
	Temperature      float32
	TopP             float32
	N                int
	PresencePenalty  float32
	FrequencyPenalty float32
	User             string
}

func WithMaxToken(token int) Option {
	return func(o *Options) {
		o.MaxToken = token
	}
}

func WithModel(model string) Option {
	return func(o *Options) {
		o.Model = model
	}
}

func WithTemperature(temperature float32) Option {
	return func(o *Options) {
		o.Temperature = temperature
	}
}

func WithTopP(topP float32) Option {
	return func(o *Options) {
		o.TopP = topP
	}
}

func WithN(n int) Option {
	return func(o *Options) {
		o.N = n
	}
}

func WithPresencePenalty(presencePenalty float32) Option {
	return func(o *Options) {
		o.PresencePenalty = presencePenalty
	}
}

func WithFrequencyPenalty(frequencyPenalty float32) Option {
	return func(o *Options) {
		o.FrequencyPenalty = frequencyPenalty
	}
}

func WithUser(user string) Option {
	return func(o *Options) {
		o.User = user
	}
}
