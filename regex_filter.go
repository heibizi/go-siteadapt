package siteadapt

import (
	"fmt"
	"regexp"
)

type regexFilter struct {
	text   string
	filter *Filter
}

func (f regexFilter) doFilter() (string, error) {
	if arg, ok := f.filter.Args.(string); ok {
		re := regexp.MustCompile(arg)
		matches := re.FindStringSubmatch(f.text)
		if matches != nil {
			return "1", nil
		} else {
			return "0", nil
		}
	} else if args, ok := f.filter.Args.([]any); ok && len(args) == 3 {
		if arg, isStr := args[0].(string); isStr {
			re := regexp.MustCompile(arg)
			matches := re.FindStringSubmatch(f.text)
			if matches != nil {
				return fmt.Sprint(args[1]), nil
			} else {
				return fmt.Sprint(args[2]), nil
			}
		}
	}
	return "", fmt.Errorf("过滤器参数格式错误: %v", f.filter.Args)
}
