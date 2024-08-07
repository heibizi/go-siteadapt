package siteadapt

import (
	"fmt"
)

type textFilter interface {
	// 过滤文本，过滤不成功请返回原始文本，交给后续 Filter 处理
	doFilter() (string, error)
}

type fieldFilterType string

const (
	fieldFilterTypeResearch    fieldFilterType = "re_search"
	fieldFilterTypeSplit       fieldFilterType = "split"
	fieldFilterTypeReplace     fieldFilterType = "replace"
	fieldFilterTypeStrip       fieldFilterType = "strip"
	fieldFilterTypeAppendLeft  fieldFilterType = "append_left"
	fieldFilterTypeQueryString fieldFilterType = "querystring"
	fieldFilterTypeRegex       fieldFilterType = "regex"
	fieldFilterTypeByteSize    fieldFilterType = "byte_size"
	fieldFilterTypeTimestamp   fieldFilterType = "timestamp"
	fieldFilterTypeEq          fieldFilterType = "eq"
	fieldFilterTypeCase        fieldFilterType = "case"
	fieldFilterTypeNotBlank    fieldFilterType = "not_blank"
	fieldFilterTypeBlank       fieldFilterType = "blank"
	fieldFilterTypeConstant    fieldFilterType = "constant"
)

func newTextFilter(text string, filter *Filter) (textFilter, error) {
	filterName := filter.Name
	if filterName == string(fieldFilterTypeResearch) {
		return researchFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeSplit) {
		return splitFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeReplace) {
		return replaceFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeStrip) {
		return stripFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeAppendLeft) {
		return appendLeftFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeQueryString) {
		return queryStringFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeRegex) {
		return regexFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeByteSize) {
		return byteSizeFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeTimestamp) {
		return timestampFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeEq) {
		return eqFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeCase) {
		return caseFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeNotBlank) {
		return notBlankFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeBlank) {
		return blankFilter{text, filter}, nil
	} else if filterName == string(fieldFilterTypeConstant) {
		return constantFilter{text, filter}, nil
	} else {
		return nil, fmt.Errorf("unknown text filter: %s", filterName)
	}
}
