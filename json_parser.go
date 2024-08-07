package siteadapt

import (
	"github.com/tidwall/gjson"
)

type jsonParser struct {
	rd RequestDefinition
}

func (p *jsonParser) read(data []byte) (gjson.Result, error) {
	return gjson.ParseBytes(data), nil
}

func (p *jsonParser) each(item gjson.Result, list List) []gjson.Result {
	return item.Get(list.Selector).Array()
}

func (p *jsonParser) get(item gjson.Result, field Field) gjson.Result {
	return item.Get(p.defaultSelector(field))
}

func (p *jsonParser) parseArray(item gjson.Result, field Field) ([]string, error) {
	var texts []string
	for _, result := range item.Get(p.defaultSelector(field)).Array() {
		text, err := p.parse(result, field)
		if err != nil {
			return nil, err
		}
		if len(text) > 0 {
			texts = append(texts, text)
		}
	}
	return texts, nil
}

func (p *jsonParser) parseString(item gjson.Result, field Field) (string, error) {
	return p.parse(item.Get(p.defaultSelector(field)), field)
}

func (p *jsonParser) parse(item gjson.Result, field Field) (string, error) {
	return filterText(item.String(), field.Filters)
}

func (p *jsonParser) defaultSelector(field Field) string {
	if field.Selector == "" {
		return field.Name
	}
	return field.Selector
}
