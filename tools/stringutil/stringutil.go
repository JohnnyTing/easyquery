package stringutil

import (
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
)

func IsEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func IsEmptySilce(str []string) bool {
	return len(str) == 0
}

func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

func StrConcat(arr ...string) string {
	return strings.Join(arr, "")
}

func StrContains(items []string, item string) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}
	return false
}

func StrInsensitiveContains(items []string, item string) (string, bool) {
	for _, v := range items {
		if strcase.ToCamel(v) == strcase.ToCamel(item) {
			return strcase.ToCamel(item), true
		}
	}
	return "", false
}

func StrInsenContains(item string, items ...string) bool {
	for _, v := range items {
		if strcase.ToCamel(v) == strcase.ToCamel(item) {
			return true
		}
	}
	return false
}

func RetrieveItem(items []string, item string) (string, bool) {
	for _, obj := range items {
		if obj == item {
			return item, true
		}
	}
	return "", false
}

func StrInsensitive(one, other string) bool {
	return strcase.ToCamel(one) == strcase.ToCamel(other)
}

func Str2Uint(str string) uint {
	result := Str2Int(str)
	return uint(result)
}

func Str2Int(str string) int {
	result, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return result
}

func Str2Int64(str string) int64 {
	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return result
}

func Str2Float64(str string) float64 {
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return result
}

func Str2Float64P(str string) *float64 {
	result := Str2Float64(str)
	return &result
}

func Str2Float32(str string) float32 {
	result, err := strconv.ParseFloat(str, 32)
	if err != nil {
		panic(err)
	}
	return float32(result)
}

func Marshal(v interface{}) []byte {
	data, err := jsoniter.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func MarshalToString(v interface{}) string {
	data, err := jsoniter.MarshalToString(v)
	if err != nil {
		panic(err)
	}
	return data
}

func Unmarshal(data []byte, v interface{}) {
	err := jsoniter.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}
