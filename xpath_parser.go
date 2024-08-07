package siteadapt

import (
	"bytes"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type xPathParser struct {
	rd RequestDefinition
}

func (p *xPathParser) read(data []byte) (*html.Node, error) {
	doc, err := htmlquery.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("解析 html 异常: %v", err)
	}
	return doc, nil
}

func (p *xPathParser) each(item *html.Node, list List) []*html.Node {
	return htmlquery.Find(item, list.Selector)
}

func (p *xPathParser) get(item *html.Node, field Field) *html.Node {
	return htmlquery.FindOne(item, field.Selector)
}

func (p *xPathParser) parseArray(item *html.Node, field Field) ([]string, error) {
	var texts []string
	for _, node := range htmlquery.Find(item, field.Selector) {
		text, err := p.parse(node, field)
		if err != nil {
			return nil, err
		}
		if len(text) > 0 {
			texts = append(texts, text)
		}
	}
	return texts, nil
}

func (p *xPathParser) parseString(item *html.Node, field Field) (string, error) {
	return p.parse(htmlquery.FindOne(item, field.Selector), field)
}

func (p *xPathParser) parse(item *html.Node, field Field) (string, error) {
	text := ""
	if item != nil {
		if field.Selection == string(fieldSelectionHtml) {
			text = htmlquery.OutputHTML(item, true)
		} else {
			text = htmlquery.InnerText(item)
		}
	}
	return filterText(text, field.Filters)
}
