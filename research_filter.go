package siteadapt

import (
	"fmt"
	"regexp"
)

type researchFilter struct {
	text   string
	filter *Filter
}

func (f researchFilter) doFilter() (string, error) {
	if args, ok := f.filter.Args.([]any); ok && len(args) == 2 {
		pattern, ok1 := args[0].(string)
		index, ok2 := args[1].(float64) // JSON numbers are unmarshalled as float64
		if ok1 && ok2 {
			re := regexp.MustCompile(pattern)
			matches := re.FindStringSubmatch(f.text)
			if matches != nil && int(index) < len(matches) {
				return matches[int(index)], nil // Skip the first element which is the whole match
			} else {
				return "", nil
			}
		}
	}
	return f.text, fmt.Errorf("过滤器参数格式错误: %v", f.filter.Args)
}
