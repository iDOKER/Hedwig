package common

import (
	"regexp"
)

// Filter 过滤器函数
// 入参：内容（string），过滤模式（string）
// 出参：过滤后的内容（string），错误（error）
func Filter(content, pattern string) (string, error) {
	// 编译正则表达式
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	// 使用正则表达式替换内容
	filteredContent := re.ReplaceAllString(content, "")

	return filteredContent, nil
}
