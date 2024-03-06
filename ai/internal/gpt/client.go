package gpt

import (
	"github.com/sashabaranov/go-openai"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	Token        string        `json:",optional"`
	TransportUrl string        `json:",optional"`
	Timeout      time.Duration `json:",default=10s"`
	MaxToken     int           `json:",default=4000"`
}

var client *openai.Client

func Setup(config Config) {
	c := openai.DefaultAzureConfig(config.Token, "https://minduck-openai.openai.azure.com/")
	Image.openaiSk = config.Token

	var transport *http.Transport

	if config.TransportUrl != "" {
		// create HTTP Transport object and set the proxy server
		// 创建一个 HTTP Transport 对象，并设置代理服务器
		proxyUrl, err := url.Parse(config.TransportUrl)
		if err != nil {
			panic(err)
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}

	// create HTTP client and set the Transport object to its Transport field
	// 创建一个 HTTP 客户端，并将 Transport 对象设置为其 Transport 字段
	c.HTTPClient.Timeout = config.Timeout
	if transport != nil {
		c.HTTPClient.Transport = transport
	}

	client = openai.NewClientWithConfig(c)
}
