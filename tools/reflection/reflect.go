package reflection

import (
	"easyquery/tools/constant"
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
)

// 返回字段值(lower_case)
func TransferFields(model interface{}, join bool) []string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Invalid {
		panic("reflection is invalid")
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("Check type error not Struct")
	}
	joinTable := strcase.ToLowerCamel(t.Name())
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := t.Field(i)
		result = extractor(field, result, join, joinTable)
	}
	return result
}

func extractor(field reflect.StructField, result []string, join bool, joinTable string) []string {
	if field.Type.Kind() == reflect.Struct && field.Anonymous {
		inField := field.Type
		for j := 0; j < inField.NumField(); j++ {
			result = extractor(inField.Field(j), result, join, joinTable)
		}
	} else {
		result = addFields(result, field, join, joinTable)
	}
	return result
}

func addFields(data []string, field reflect.StructField, join bool, joinTable string) []string {
	if join {
		// company_j_name
		tagName := fmt.Sprintf("%s%s%s", joinTable, constant.Delimiter, strcase.ToLowerCamel(field.Name))
		data = append(data, tagName)
	} else {
		// 前后端字段为小写驼峰命名
		tagName := strcase.ToLowerCamel(field.Name)
		data = append(data, tagName)
	}
	return data
}

func InvokeMethod(model interface{}, method string, params ...interface{}) (bool, interface{}) {
	v := reflect.ValueOf(model)
	m := v.MethodByName(method)
	if m.Kind() == reflect.Invalid {
		return false, ""
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return true, m.Call(in)[0].Interface()
}
