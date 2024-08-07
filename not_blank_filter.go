package siteadapt

import (
	"fmt"
)

type notBlankFilter struct {
	text   string
	filter *Filter
}

func (f notBlankFilter) doFilter() (string, error) {
	if f.filter.Args == nil {
		if len(f.text) > 0 {
			return "1", nil
		} else {
			return "0", nil
		}
	} else if arg, ok := f.filter.Args.(string); ok {
		if len(f.text) > 0 {
			return arg, nil
		} else {
			return f.text, nil
		}
	} else if args, ok := f.filter.Args.([]any); ok && len(args) == 2 {
		if len(f.text) > 0 {
			return fmt.Sprint(args[0]), nil
		} else {
			return fmt.Sprint(args[1]), nil
		}
	}
	return f.text, fmt.Errorf("过滤器参数格式错误: %v", f.filter)
}
