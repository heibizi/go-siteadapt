package siteadapt

import (
	"encoding/json"
	"fmt"
)

type ConfigReader struct {
	data []byte
}

func NewConfigReader(data []byte) *ConfigReader {
	return &ConfigReader{
		data: data,
	}
}

func (c *ConfigReader) Read() (*Config, error) {
	var data interface{}
	err := json.Unmarshal(c.data, &data)
	if err != nil {
		return nil, fmt.Errorf("解析配置异常: %v", err)
	}
	var config Config
	if err := WeakDecode(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置异常: %v", err)
	}
	c.initialConfig(&config)
	return &config, nil
}

func (c *ConfigReader) initialConfig(config *Config) {
	for name, rd := range config.RequestDefinitions {
		// 先引用公共字段过来
		if len(rd.FieldsRef) > 0 {
			// 字段引用
			rd.Fields = config.CommonFields[rd.FieldsRef]
			config.RequestDefinitions[name] = rd
		}
		// 再做相关设置
		for name, field := range rd.Fields {
			c.setupField(&field, name, config)
			rd.Fields[name] = field
		}
	}
}

func (c *ConfigReader) setupField(field *Field, name string, config *Config) {
	field.Name = name
	if field.Any != nil {
		// 把外层 filters 复制到内层
		for i := range field.Any {
			anyFieldPtr := &field.Any[i]
			anyFieldPtr.Name = name
			for _, filter := range field.Filters {
				anyFieldPtr.Filters = append(anyFieldPtr.Filters, filter)
			}
		}
		// 清空外层的 filters
		field.Filters = nil
	}
	if field.FieldsRef != "" {
		field.Fields = config.CommonFields[field.FieldsRef]
	}
	if field.Fields != nil {
		for name, subField := range field.Fields {
			c.setupField(&subField, name, config)
			field.Fields[name] = subField
		}
	}
}
