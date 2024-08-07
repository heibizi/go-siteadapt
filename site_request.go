package siteadapt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type (
	siteRequest struct {
		sc     Config
		params params

		requestUrl string
	}
	params struct {
		rd       RequestDefinition // 自定义请求
		headers  map[string]string // 请求头
		params   url.Values        // url 请求参数
		formData url.Values        // form-data 请求参数
		body     map[string]any    // 请求体
		env      map[string]string // 环境变量
		ua       string
		cookie   string
	}
)

func newSiteRequest(sc Config, params params) *siteRequest {
	return &siteRequest{
		sc:     sc,
		params: params,
	}
}

func (cr *siteRequest) newRequest(rd RequestDefinition) (*http.Request, error) {
	// 占位符替换
	rd.Params = cr.replaceUrlValues(rd.Params)
	for k, vs := range cr.params.params {
		for _, v := range vs {
			rd.Params.Add(k, v)
		}
	}
	rd.FormData = cr.replaceUrlValues(rd.FormData)
	if cr.params.formData != nil {
		for k, vs := range cr.params.formData {
			for _, v := range vs {
				rd.FormData.Add(k, v)
			}
		}
	}
	requestUrl := cr.newRequestUrl(rd)
	if len(rd.Params) > 0 {
		requestUrl = fmt.Sprintf("%s?%s", requestUrl, rd.Params.Encode())
	}
	var requestBody io.Reader
	if len(rd.FormData) > 0 {
		// 占位符替换
		requestBody = bytes.NewBufferString(rd.FormData.Encode())
	} else if cr.params.body != nil {
		// 将数据序列化为 JSON
		jsonData, err := json.Marshal(cr.params.body)
		if err != nil {
			return nil, cr.newError("%s 请求体转json异常: %v, 异常: %v", requestUrl, cr.params.body, err)
		}
		requestBody = bytes.NewReader(jsonData)
	}
	requestUrl = cr.replacePlaceholders(requestUrl)
	// todo 后续记得增加代理，如果站点设置了代理的话
	req, err := http.NewRequest(rd.Method, requestUrl, requestBody)
	if err != nil {
		return nil, cr.newError("请求站点异常: %v", err)
	}

	// header 的优先级：自定义请求头 > 站点配置请求头 > 默认请求头
	// 默认请求头
	req.Header.Set("User-Agent", cr.params.ua)
	if len(cr.params.cookie) > 0 {
		req.Header.Set("Cookie", cr.params.cookie)
	}
	// 站点配置请求头
	if rd.Headers != nil {
		for k, v := range rd.Headers {
			v = cr.replacePlaceholders(v)
			req.Header.Set(k, v)
		}
	}
	// 自定义请求头
	for k, v := range cr.params.headers {
		req.Header.Set(k, v)
	}
	// 校验必填请求头
	if len(rd.RequiredHeaders) > 0 {
		for _, requiredHeader := range rd.RequiredHeaders {
			if len(req.Header.Get(requiredHeader)) == 0 {
				return nil, cr.newError("请求头 %s 未设置", requiredHeader)
			}
		}
	}
	cr.requestUrl = requestUrl
	return req, nil
}

func (cr *siteRequest) httpRequest() (*http.Response, error) {
	client := &http.Client{}
	request, err := cr.newRequest(cr.params.rd)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, cr.newError("请求站点异常: %v", err)
	}
	return resp, nil
}

func (cr *siteRequest) newError(format string, v ...any) error {
	return fmt.Errorf("站点(%s)%s", cr.sc.Name, fmt.Sprintf(format, v...))
}

func (cr *siteRequest) replacePlaceholders(str string) string {
	for placeholder, replacement := range cr.params.env {
		str = strings.ReplaceAll(str, "{"+placeholder+"}", replacement)
	}
	return str
}
func (cr *siteRequest) replaceUrlValues(values url.Values) url.Values {
	newValues := url.Values{}
	for key, vals := range values {
		for _, val := range vals {
			newValues.Add(key, cr.replacePlaceholders(val))
		}
	}
	return newValues
}

func (cr *siteRequest) newRequestUrl(rd RequestDefinition) string {
	if IsValidHttpUrl(rd.Path) {
		return rd.Path
	} else {
		baseUrl := ""
		if rd.UseApi {
			baseUrl = cr.sc.Api
		} else {
			baseUrl = cr.sc.Domain
		}
		return baseUrl + rd.Path
	}
}
