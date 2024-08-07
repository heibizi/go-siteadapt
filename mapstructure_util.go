package siteadapt

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// 自定义 DecodeHook 函数
func decodeHookFunc(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	// 字符串转其他类型
	if f.Kind() == reflect.String {
		str := data.(string)
		if t.Kind() == reflect.Int {
			return ParseInt(str), nil
		}
		if t.Kind() == reflect.Int64 {
			return ParseInt64(str), nil
		}
		if t.Kind() == reflect.Float64 {
			return ParseFloat64(str), nil
		}
		if t.Kind() == reflect.Bool {
			return ParseBool(str), nil
		}
	}
	// float64 转其他类型，json 里 int 类型也是 float64
	if f.Kind() == reflect.Float64 {
		number := data.(float64)
		if t.Kind() == reflect.Int {
			return int(number), nil
		}
		if t.Kind() == reflect.Int64 {
			return int64(number), nil
		}
		if t.Kind() == reflect.String {
			return fmt.Sprint(number), nil
		}
	}
	return data, nil
}

// WeakDecode 自定义 weakDecode 函数，转换 string 为其他类型，取决于 struct 定义的类型。
func WeakDecode(input, output interface{}) error {
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			decodeHookFunc,
		),
		Result: output,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return err
	}
	if err := decoder.Decode(input); err != nil {
		return err
	}
	return decoder.Decode(input)
}
