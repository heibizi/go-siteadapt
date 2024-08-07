package siteadapt

import (
	"fmt"
)

type timestampFilter struct {
	text   string
	filter *Filter
}

func (f timestampFilter) doFilter() (string, error) {
	return fmt.Sprintf("%d", GetTimeStamp(f.text)), nil
}
