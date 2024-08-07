package siteadapt

import (
	"fmt"
)

type constantFilter struct {
	text   string
	filter *Filter
}

func (f constantFilter) doFilter() (string, error) {
	if arg, ok := f.filter.Args.(string); ok {
		return arg, nil
	}
	return f.text, fmt.Errorf("过滤器参数格式错误: %v", f.filter.Args)
}
