package gpt

import "regexp"

func GetMarkdown(s string) string {
	// 从 s 字符串中，使用正则表达式获取 ```*\n 的位置，例如 ```\n ```text\n ```markdown\n
	re := regexp.MustCompile("```.*\n")
	startIndex := re.FindStringIndex(s)
	if startIndex == nil || len(startIndex) != 2 {
		return s
	}

	s = s[startIndex[1]:]
	re2 := regexp.MustCompile("\n```")
	endIndex := re2.FindStringIndex(s)
	if endIndex == nil || len(endIndex) != 2 {
		return s
	}
	return s[:endIndex[0]]
}
