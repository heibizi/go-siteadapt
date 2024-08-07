package siteadapt

import (
	"fmt"
)

type caseFilter struct {
	text   string
	filter *Filter
}

func (f caseFilter) doFilter() (string, error) {
	if slice, ok := f.filter.Args.(map[string]any); ok {
		if text, exists := slice[f.text]; exists {
			return fmt.Sprint(text), nil
		}
		if defaultValue, exists := slice["*"]; exists {
			return fmt.Sprint(defaultValue), nil
		}
		return f.text, nil
	}
	return f.text, fmt.Errorf("过滤器参数格式错误: %v", f.filter.Args)
}
