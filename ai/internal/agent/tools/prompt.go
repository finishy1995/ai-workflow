package tools

import (
	"fmt"
	"strings"
)

func MergeParam(priorityDesc ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range priorityDesc {
		if m != nil {
			for key, value := range m {
				if _, ok := result[key]; !ok {
					result[key] = value
				}
			}
		}
	}
	return result
}

func FillParam(param map[string]interface{}, prompts []string) []string {
	results := make([]string, 0, len(prompts))
	for _, prompt := range prompts {
		result := prompt
		for k, v := range param {
			result = strings.ReplaceAll(result, fmt.Sprintf("{%s}", k), fmt.Sprintf("%v", v))
		}
		results = append(results, result)
	}
	return results

}

func FillParamForString(param map[string]interface{}, prompt string) string {
	result := prompt
	for k, v := range param {
		result = strings.ReplaceAll(result, fmt.Sprintf("{%s}", k), fmt.Sprintf("%v", v))
	}
	return result
}
