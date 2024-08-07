package siteadapt

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"io"
)

func NewSiteAdaptor(sc Config) *SiteAdaptor {
	return &SiteAdaptor{
		sc: sc,
	}
}

type (
	SiteAdaptor struct {
		sc Config
	}

	// RequestSiteResult 站点请求结果
	result struct {
		List     []map[string]any // 列表数据
		Raw      []byte           // 原始数据
		Data     map[string]any   // 键值对数据
		NextPage string           // 下一页
		RequestInfo
	}
	RequestInfo struct {
		Domain     string // 域名
		RequestUrl string // 请求地址
		StatusCode int    // 状态码
	}

	ListResult struct {
		NextPage string
		RequestInfo
	}

	DataResult struct {
		RequestInfo
	}

	RawResult struct {
		Data []byte
		RequestInfo
	}

	ListFunc func(result ListResult)
	DataFunc func(result DataResult)
	RawFunc  func(result RawResult)
)

// List 获取列表数据
func (sa *SiteAdaptor) List(requestSiteParams RequestSiteParams, output any, fn ListFunc) error {
	r, err := sa.requestSite(requestSiteParams)
	if err != nil {
		return fmt.Errorf("请求异常: %v", err)
	}
	err = WeakDecode(r.List, output)
	if err != nil {
		return fmt.Errorf("解析异常: %v", err)
	}
	if fn != nil {
		fn(ListResult{
			NextPage:    r.NextPage,
			RequestInfo: r.RequestInfo,
		})
	}
	return nil
}

// Data 获取对象数据
func (sa *SiteAdaptor) Data(requestSiteParams RequestSiteParams, output any, fn DataFunc) error {
	r, err := sa.requestSite(requestSiteParams)
	if err != nil {
		return fmt.Errorf("请求异常: %v", err)
	}
	err = WeakDecode(r.Data, output)
	if err != nil {
		return fmt.Errorf("解析异常: %v", err)
	}
	if fn != nil {
		fn(DataResult{
			r.RequestInfo,
		})
	}
	return nil
}

// Raw 获取原始数据
func (sa *SiteAdaptor) Raw(requestSiteParams RequestSiteParams, fn RawFunc) error {
	r, err := sa.requestSite(requestSiteParams)
	if err != nil {
		return fmt.Errorf("请求异常: %v", err)
	}
	if fn != nil {
		fn(RawResult{
			Data:        r.Raw,
			RequestInfo: r.RequestInfo,
		})
	}
	return nil
}

// Json 直接 JSON 转 Struct
func (sa *SiteAdaptor) Json(requestSiteParams RequestSiteParams, output interface{}) error {
	r, err := sa.requestSite(requestSiteParams)
	if err != nil {
		return fmt.Errorf("请求异常: %v", err)
	}
	var input map[string]interface{}
	err = json.Unmarshal(r.Raw, &input)
	if err != nil {
		return err
	}
	err = mapstructure.WeakDecode(input, &output)
	if err != nil {
		return fmt.Errorf("解析异常: %v", err)
	}
	return nil
}

// 请求站点
func (sa *SiteAdaptor) requestSite(requestSiteParams RequestSiteParams) (result, error) {
	reqId := requestSiteParams.ReqId
	rd, err := sa.getRd(reqId, requestSiteParams)
	if err != nil {
		return result{}, err
	}
	sc := sa.sc
	if requestSiteParams.Domain != "" {
		sc.Domain = requestSiteParams.Domain
	}
	if requestSiteParams.Api != "" {
		sc.Api = requestSiteParams.Api
	}
	if requestSiteParams.Path != "" {
		rd.Path = requestSiteParams.Path
	}
	sr := newSiteRequest(sc, params{
		rd:       rd,
		headers:  requestSiteParams.Headers,
		params:   requestSiteParams.Params,
		formData: requestSiteParams.FormData,
		body:     requestSiteParams.Body,
		env:      requestSiteParams.Env,
		ua:       requestSiteParams.UA,
		cookie:   requestSiteParams.Cookie,
	})
	resp, err := sr.httpRequest()
	if err != nil {
		return result{}, err
	}
	defer resp.Body.Close()
	success := false
	if rd.SuccessStatusCodes == nil {
		if resp.StatusCode == 200 {
			success = true
		}
	} else {
		for _, successStatusCode := range rd.SuccessStatusCodes {
			if resp.StatusCode == successStatusCode {
				success = true
				break
			}
		}
	}
	requestUrl := sr.requestUrl
	if !success {
		return result{}, sa.newError("%s请求失败, 状态码为: %v, 异常: %v", requestUrl, resp.StatusCode, err)
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return result{}, sa.newError("请求异常: %v", err)
	}
	// 解析响应
	parsedData, err := newParserHelper(responseBody, sa.sc, rd).parse()
	if err != nil {
		return result{}, sa.newError("解析响应异常: %v", err)
	}
	requestInfo := RequestInfo{
		Domain:     sc.Domain,
		RequestUrl: requestUrl,
		StatusCode: resp.StatusCode,
	}
	result := result{
		RequestInfo: requestInfo,
	}
	nextPage := parsedData[string(fieldNameNextPage)]
	if nextPage != nil {
		result.NextPage = nextPage.(string)
	}
	if parsedData[string(fieldNameList)] != nil {
		result.List = parsedData[string(fieldNameList)].([]map[string]any)
	} else if parsedData[string(fieldNameRaw)] != nil {
		result.Raw = parsedData[string(fieldNameRaw)].([]byte)
	} else {
		result.Data = parsedData
	}
	return result, nil
}

func (sa *SiteAdaptor) newError(format string, v ...any) error {
	return fmt.Errorf("站点(%s)%s", sa.sc.Name, fmt.Sprintf(format, v...))
}

func (sa *SiteAdaptor) getRd(reqId string, params RequestSiteParams) (RequestDefinition, error) {
	if params.Rd != nil {
		return *params.Rd, nil
	}
	scRd, exists := sa.sc.RequestDefinitions[reqId]
	if exists {
		return scRd, nil
	}
	return RequestDefinition{}, sa.newError("未适配该请求: %s", reqId)
}
