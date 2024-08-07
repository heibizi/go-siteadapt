package siteadapt

import (
	"net/url"
)

type (
	// Config 站点配置
	Config struct {
		ID                 string                       `mapstructure:"id"`            // 站点 ID
		Name               string                       `mapstructure:"name"`          // 站点名称
		Domain             string                       `mapstructure:"domain"`        // 域名
		Api                string                       `mapstructure:"api"`           // api 地址
		RequestDefinitions map[string]RequestDefinition `mapstructure:"requests"`      // 请求定义
		CommonFields       map[string]map[string]Field  `mapstructure:"common_fields"` // 公共字段定义
	}
	// Field 站点请求字段定义
	Field struct {
		Name           string           // 不需要填写
		List           *List            `mapstructure:"list,optional"` // 列表数据配置
		Selector       string           `mapstructure:"selector"`      // 选择器表达式，不同解析类型有不同的表达式
		Selection      string           `mapstructure:"selection"`     // 选择器文本类型，见：FieldSelection
		Attribute      string           `mapstructure:"attribute"`     // 用于 HTML 解析
		Filters        []Filter         `mapstructure:"filters"`       // 过滤器，链式处理
		Parent         bool             `mapstructure:"parent"`        // 是否选择父元素
		ChildrenRemove string           `mapstructure:"remove"`        // 用于 HTML 解析，移除指定元素的子元素
		Any            []Field          `mapstructure:"any"`           // 数组任意匹配，直到匹配成功
		Array          bool             `mapstructure:"array"`         // 是否为数组
		Fields         map[string]Field `mapstructure:"fields"`        // 字段嵌套，意味着可以多层解析
		FieldsRef      string           `mapstructure:"fields_ref"`    // 字段嵌套引用 CommonFields
	}
	// Filter 站点请求字段过滤器
	Filter struct {
		Name string `mapstructure:"name"`
		Args any    `mapstructure:"args"`
	}
	List struct {
		Selector string `mapstructure:"selector"`
		NextPage Field  `mapstructure:"next_page"`
	}
	// RequestDefinition 站点请求定义 RequestDefinition
	RequestDefinition struct {
		Parser             string            `mapstructure:"parser"`               // 解析器，见：requestParser
		Method             string            `mapstructure:"method"`               // 请求方法
		Path               string            `mapstructure:"path"`                 // 相对路径
		UseApi             bool              `mapstructure:"use_api"`              // 使用 api
		Chrome             bool              `mapstructure:"chrome"`               // 使用 chrome
		Headers            map[string]string `mapstructure:"headers"`              // 请求头
		RequiredHeaders    []string          `mapstructure:"required_headers"`     // 必填请求头
		Params             url.Values        `mapstructure:"params"`               // 请求参数
		FormData           url.Values        `mapstructure:"form_data"`            // 表单数据
		SuccessStatusCodes []int             `mapstructure:"success_status_codes"` // 成功状态码
		List               *List             `mapstructure:"list,optional"`        // 列表数据配置
		Fields             map[string]Field  `mapstructure:"fields"`               // 字段定义
		FieldsRef          string            `mapstructure:"fields_ref"`           // 引用 CommonFields
		DisabledExtends    struct {
			Path               bool `mapstructure:"path"`
			Headers            bool `mapstructure:"headers"`
			RequiredHeaders    bool `mapstructure:"required_headers"`
			Params             bool `mapstructure:"params"`
			FormData           bool `mapstructure:"form_data"`
			SuccessStatusCodes bool `mapstructure:"success_status_codes"`
			List               bool `mapstructure:"list"`
			Fields             bool `mapstructure:"fields"`
			Field              bool `mapstructure:"field"`
		} `mapstructure:"disabled_extends"` // 禁用继承
	}
)

type (
	// RequestSiteParams 站点请求参数
	RequestSiteParams struct {
		ReqId    string             // 请求 id
		Rd       *RequestDefinition // 自定义 rd
		Domain   string             // 自定义 domain
		Api      string             // 自定义 api
		Path     string             // 自定义 path
		Headers  map[string]string  // 请求头
		Params   url.Values         // url 请求参数
		FormData url.Values         // form-data 请求参数
		Body     map[string]any     // 请求体
		Env      map[string]string  // 环境变量
		UA       string
		Cookie   string
	}
)
