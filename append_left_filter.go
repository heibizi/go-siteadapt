package siteadapt

import (
	"fmt"
)

type appendLeftFilter struct {
	text   string
	filter *Filter
}

func (f appendLeftFilter) doFilter() (string, error) {
	prefix, ok := f.filter.Args.(string)
	if ok {
		return prefix + f.text, nil
	}
	return f.text, fmt.Errorf("过滤器参数格式错误: %v", f.filter)
}
