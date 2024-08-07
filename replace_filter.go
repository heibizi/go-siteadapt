package siteadapt

import (
	"fmt"
	"strings"
)

type replaceFilter struct {
	text   string
	filter *Filter
}

func (f replaceFilter) doFilter() (string, error) {
	if args, ok := f.filter.Args.([]any); ok && len(args) == 2 {
		old, ok1 := args[0].(string)
		newText, ok2 := args[1].(string)
		if ok1 && ok2 {
			return strings.ReplaceAll(f.text, old, newText), nil
		}
	}
	return f.text, fmt.Errorf("过滤器参数格式错误: %v", f.filter.Args)
}
