package siteadapt

import (
	"fmt"
	"strings"
)

type queryStringFilter struct {
	text   string
	filter *Filter
}

func (f queryStringFilter) doFilter() (string, error) {
	if s, ok := f.filter.Args.(string); ok {
		index := strings.Index(f.text, s)
		if index >= 0 {
			return f.text[index:], nil
		}
		return f.text, nil
	}
	return f.text, fmt.Errorf("过滤器参数格式错误: %v", f.filter.Args)
}
