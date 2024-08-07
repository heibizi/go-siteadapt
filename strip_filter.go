package siteadapt

import (
	"strings"
)

type stripFilter struct {
	text   string
	filter *Filter
}

func (f stripFilter) doFilter() (string, error) {
	return strings.TrimSpace(f.text), nil
}
