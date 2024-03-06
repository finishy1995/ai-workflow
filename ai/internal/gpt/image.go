package gpt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// todo refactor
type image struct {
	client             *http.Client
	openaiSk           string
	stableDiffusionReq string
	stableDiffusionSK  string
	dalle3Req          string
}

var Image = &image{
	client: &http.Client{
		Timeout: 120 * time.Second,
	},
	openaiSk: "",
	stableDiffusionReq: `{
    "key": "3GshGJ472T3wnia1fHf7uYcuW73XvlTWDTusXkW4CRAxRGqRBqCfVRTU73Mm",
    "model_id": "midjourney",
    "prompt": "{{%s}}",
    "width": "512",
    "height": "512",
    "samples": "2",
    "num_inference_steps": "30",
    "seed": null,
    "guidance_scale": 7.5,
    "webhook": null,
    "track_id": null
}`,
	stableDiffusionSK: `{
    "key": "3GshGJ472T3wnia1fHf7uYcuW73XvlTWDTusXkW4CRAxRGqRBqCfVRTU73Mm"
}`,

	dalle3Req: `{
    "model": "dall-e-3",
    "prompt": "{{%s}}",
    "n": 1,
    "size": "1024x1024",
    "style": "{{%t}}",
    "quality": "hd"
}`,
}

func (img *image) CreateImage(answer string) string {
	var sg sync.WaitGroup
	sg.Add(3)
	var r1, r2, r3s string

	go func() {
		defer sg.Done()
		r1 = img.getDallE3Result(answer, openai.CreateImageStyleVivid)
	}()

	go func() {
		defer sg.Done()
		r2 = img.getDallE3Result(answer, openai.CreateImageStyleVivid)
	}()

	go func() {
		defer sg.Done()
		r3s = img.getStableDiffusionResult(answer)
	}()

	sg.Wait()
	return fmt.Sprintf("%s;%s;%s", r1, r2, r3s)
}

func (img *image) getStableDiffusionResult(answer string) (resultStableDiffusion string) {

	payload := strings.NewReader(strings.Replace(img.stableDiffusionReq, "{{%s}}", answer, -1))
	req, err := http.NewRequest("POST", "https://stablediffusionapi.com/api/v3/dreambooth", payload)
	if err != nil {
		logx.Errorf("image getActualAnswer failed, error: %s", err.Error())
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := img.client.Do(req)
	if err != nil {
		logx.Errorf("image getActualAnswer failed, error: %s", err.Error())
		return
	}
	//{
	//	"status": "success",
	//	"generationTime": 10.171663284301758,
	//	"id": 32401,
	//	"output": [
	//		"https://stable-diffusion-api.de-fra1.upcloudobjects.com/generations/018704228461667621441.png"
	//	]
	//}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("image getActualAnswer failed, error: %s", err.Error())
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)

	if err != nil {
		logx.Errorf("image getActualAnswer json unmarshal failed, error: %s", err.Error())
		return
	}

	fmt.Println(result)
	if result["status"] != nil {
		if result["status"] == "success" {
			resultArr := result["output"].([]interface{})

			for index, value := range resultArr {
				if index > 3 {
					break
				}
				resultStableDiffusion += value.(string) + ";"
			}
		} else if result["fetch_result"] != nil {
			url := result["fetch_result"].(string)
			count := 0
			for {
				if count > 20 {
					return
				}
				time.Sleep(time.Second * 5)

				payload = strings.NewReader(img.stableDiffusionSK)
				req, err = http.NewRequest("POST", url, payload)
				if err != nil {
					logx.Errorf("image getActualAnswer failed, error: %s", err.Error())
					return
				}

				req.Header.Add("Content-Type", "application/json")
				resp, err := img.client.Do(req)
				if err != nil {
					logx.Errorf("image getActualAnswer failed, error: %s", err.Error())
					return
				}

				defer resp.Body.Close()

				body, err = io.ReadAll(resp.Body)
				if err != nil {
					logx.Errorf("image getActualAnswer failed, error: %s", err.Error())
					return
				}

				var result map[string]interface{}
				err = json.Unmarshal(body, &result)

				if err != nil {
					logx.Errorf("image getActualAnswer json unmarshal failed, error: %s", err.Error())
					return
				}

				fmt.Println(result)
				if result["status"] != nil {
					if result["status"] == "success" {
						resultArr := result["output"].([]interface{})

						for index, value := range resultArr {
							if index > 3 {
								break
							}
							resultStableDiffusion += value.(string) + ";"
							return
						}
					}
				}
				count++
			}
		}
	}
	return
}

func (img *image) getDallE3Result(prompt, style string) string {
	response, err := client.CreateImage(context.Background(), openai.ImageRequest{
		Prompt:         prompt,
		Model:          openai.CreateImageModelDallE3,
		N:              1, // dalle3只能生成1个
		Quality:        "hd",
		Size:           "1024x1024",
		Style:          style,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		//User:           "",
	})
	if err != nil || len(response.Data) == 0 {
		return ""
	}
	return response.Data[0].URL
}
