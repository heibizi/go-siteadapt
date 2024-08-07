package siteadapt

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type htmlParser struct {
	rd RequestDefinition
}

func (p *htmlParser) read(data []byte) (*goquery.Selection, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("解析 html 异常: %v", err)
	}
	return doc.Selection, nil
}

func (p *htmlParser) each(item *goquery.Selection, list List) []*goquery.Selection {
	var selections []*goquery.Selection
	item.Find(list.Selector).Each(func(i int, selection *goquery.Selection) {
		selections = append(selections, selection)
	})
	return selections
}

func (p *htmlParser) get(item *goquery.Selection, field Field) *goquery.Selection {
	return item.Find(field.Selector)
}

func (p *htmlParser) parseArray(item *goquery.Selection, field Field) ([]string, error) {
	var texts []string
	list := item.Find(field.Selector)
	for i := range list.Size() {
		selection := list.Eq(i)
		text, err := p.parse(selection, field)
		if err != nil {
			return nil, err
		}
		if len(text) > 0 {
			texts = append(texts, text)
		}
	}
	return texts, nil
}

func (p *htmlParser) parseString(item *goquery.Selection, field Field) (string, error) {
	return p.parse(item.Find(field.Selector), field)
}

func (p *htmlParser) parse(item *goquery.Selection, field Field) (string, error) {
	text := ""
	if len(field.ChildrenRemove) > 0 {
		cloned := item.Clone()
		for _, selector := range strings.Split(field.ChildrenRemove, ",") {
			cloned.Find(selector).Remove()
		}
		item = cloned
	}
	if field.Parent {
		item = item.Parent()
	}
	if len(field.Attribute) > 0 {
		if val, exists := item.Attr(field.Attribute); exists {
			text = val
		}
	} else {
		if field.Selection == string(fieldSelectionHtml) {
			html, err := item.Html()
			if err != nil {
				return text, fmt.Errorf("解析 HTML 异常: %v", err)
			}
			text = html
		} else {
			text = item.Text()
		}
	}
	return filterText(text, field.Filters)
}
