package siteadapt

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
	"strings"
)

// Parser 解析器
type parser[T any] interface {
	// 读取数据
	read(data []byte) (T, error)
	// 遍历
	each(item T, list List) []T
	// 获取
	get(item T, field Field) T
	// 解析数组
	parseArray(item T, field Field) ([]string, error)
	// 解析对象
	parseString(item T, field Field) (string, error)
}

type requestParserType string

const (
	requestParserTypeCssSelector requestParserType = "CssSelector"
	requestParserTypeJsonPath    requestParserType = "JsonPath"
	requestParserTypeXPath       requestParserType = "XPath"
	requestParserTypeNone        requestParserType = "None" // 不解析
)

type fieldSelection string

const (
	fieldSelectionText fieldSelection = "text"
	fieldSelectionHtml fieldSelection = "html"
)

type fieldName string

const (
	fieldNameList     fieldName = "list"
	fieldNameRaw      fieldName = "raw"
	fieldNameNextPage fieldName = "next_page"
)

type parserHelper struct {
	data []byte
	sc   Config
	rd   RequestDefinition
}

func newParserHelper(data []byte, sc Config, rd RequestDefinition) *parserHelper {
	return &parserHelper{
		data: data,
		sc:   sc,
		rd:   rd,
	}
}

// parse 函数
func (p *parserHelper) parse() (map[string]any, error) {
	// 没有指定解析器，直接返回原始数据
	if p.rd.Parser == string(requestParserTypeNone) {
		parsedData := make(map[string]any)
		parsedData[string(fieldNameRaw)] = p.data
		return parsedData, nil
	}
	// 根据解析器类型创建解析器
	parser, err := p.newParser(p.rd)
	if err != nil {
		return nil, err
	}
	switch theParser := parser.(type) {
	case htmlParser:
		return doParse[*goquery.Selection](&theParser, p.data, p.rd)
	case xPathParser:
		return doParse[*html.Node](&theParser, p.data, p.rd)
	case jsonParser:
		return doParse[gjson.Result](&theParser, p.data, p.rd)
	default:
		return nil, fmt.Errorf("不支持的解析器类型: %s", p.rd.Parser)
	}
}

func (p *parserHelper) newParser(rd RequestDefinition) (interface{}, error) {
	if rd.Parser == string(requestParserTypeCssSelector) {
		return htmlParser{rd: rd}, nil
	} else if rd.Parser == string(requestParserTypeXPath) {
		return xPathParser{rd: rd}, nil
	} else if rd.Parser == string(requestParserTypeJsonPath) {
		return jsonParser{rd: rd}, nil
	} else {
		return nil, fmt.Errorf("不支持的解析器类型: %s", rd.Parser)
	}
}

func doParse[T any](parser parser[T], data []byte, rd RequestDefinition) (map[string]any, error) {
	doc, err := parser.read(data)
	if err != nil {
		return nil, err
	}
	return parse(parser, doc, rd.List, rd.Fields)
}

func parse[T any](parser parser[T], doc T, list *List, fields map[string]Field) (map[string]any, error) {
	parsedData := make(map[string]any)
	if list != nil {
		var arr []map[string]any
		for _, item := range parser.each(doc, *list) {
			rowData := make(map[string]any)
			for name, field := range fields {
				v, err := parseArrayOrString(parser, item, field)
				if err != nil {
					return nil, err
				}
				rowData[name] = v
			}
			arr = append(arr, rowData)
		}
		parsedData[string(fieldNameList)] = arr
		if len(list.NextPage.Selector) > 0 {
			v, err := parseArrayOrString(parser, doc, list.NextPage)
			if err != nil {
				return nil, err
			}
			parsedData[string(fieldNameNextPage)] = v
		}
	} else {
		for name, field := range fields {
			v, err := parseArrayOrString(parser, doc, field)
			if err != nil {
				return nil, err
			}
			parsedData[name] = v
		}
	}
	return parsedData, nil
}

func parseArrayOrString[T any](parser parser[T], item T, field Field) (any, error) {
	if field.Any != nil {
		for _, anyField := range field.Any {
			if field.Array {
				texts, err := parser.parseArray(item, anyField)
				if err != nil {
					return nil, err
				}
				if len(texts) > 0 {
					return texts, nil
				}
			} else {
				text, err := parser.parseString(item, anyField)
				if err != nil {
					return nil, err
				}
				if len(text) > 0 {
					return text, nil
				}
			}
		}
		return nil, nil
	} else if field.Fields != nil {
		// 字段嵌套
		if field.List != nil {
			m, err := parse(parser, item, field.List, field.Fields)
			if err != nil {
				return nil, err
			}
			return m[string(fieldNameList)], nil
		} else {
			return parse(parser, parser.get(item, field), nil, field.Fields)
		}
	} else if field.Array {
		return parser.parseArray(item, field)
	} else {
		return parser.parseString(item, field)
	}
}

// filterText 文本过滤
func filterText(text string, filters []Filter) (string, error) {
	for _, f := range filters {
		tf, err := newTextFilter(text, &f)
		if err != nil {
			return "", err
		}
		text, err = tf.doFilter()
		if err != nil {
			return "", err
		}
	}
	return strings.TrimSpace(text), nil
}
