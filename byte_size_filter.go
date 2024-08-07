package siteadapt

import (
	"fmt"
)

type byteSizeFilter struct {
	text   string
	filter *Filter
}

func (f byteSizeFilter) doFilter() (string, error) {
	return fmt.Sprintf("%d", StrToByteSize(f.text)), nil
}
