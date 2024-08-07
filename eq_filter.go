package siteadapt

import (
	"fmt"
)

type eqFilter struct {
	text   string
	filter *Filter
}

func (f eqFilter) doFilter() (string, error) {
	if arg, ok := f.filter.Args.(string); ok {
		if f.text == arg {
			return "1", nil
		} else {
			return "0", nil
		}
	} else if args, ok := f.filter.Args.([]any); ok && len(args) == 3 {
		if arg, isStr := args[0].(string); isStr && f.text == arg {
			return fmt.Sprint(args[1]), nil
		} else {
			return fmt.Sprint(args[2]), nil
		}
	}
	return f.text, fmt.Errorf("过滤器参数格式错误: %v", f.filter.Args)
}
