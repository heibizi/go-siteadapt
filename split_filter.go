package siteadapt

import (
	"fmt"
	"strings"
)

type splitFilter struct {
	text   string
	filter *Filter
}

func (f splitFilter) doFilter() (string, error) {
	if args, ok := f.filter.Args.([]any); ok && len(args) == 2 {
		delimiter, ok1 := args[0].(string)
		index, ok2 := args[1].(float64)
		if ok1 && ok2 {
			parts := strings.Split(f.text, delimiter)
			if int(index) == -1 {
				return parts[len(parts)-1], nil
			} else if int(index) < len(parts) {
				return parts[int(index)], nil
			} else {
				return f.text, fmt.Errorf("过滤器参数 index 超出分组数量: %v", f.filter.Args)
			}
		}
	}
	return f.text, fmt.Errorf("过滤器参数格式错误: %v", f.filter.Args)
}
