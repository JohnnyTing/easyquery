package easyquery

import (
	"easyquery/tools/constant"
	"easyquery/tools/reflection"
	"easyquery/tools/stringutil"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/gin-gonic/gin"
)

type QueryField struct {
	Name      string
	Type      QueryFieldType
	Operation string
	Value     interface{}
	Join      bool
	JoinTable string
}

type QueryFieldType string

const (
	String      QueryFieldType = constant.String
	Array       QueryFieldType = constant.Array
	NullOrExist QueryFieldType = constant.NullOrExist
	OrInEq      QueryFieldType = constant.OrInEq
	OrInLike    QueryFieldType = constant.OrInLike
	OrOutEq     QueryFieldType = constant.OrOutEq
	OrOutLike   QueryFieldType = constant.OrOutLike
	Order       QueryFieldType = constant.Order
	Association QueryFieldType = constant.Association
	Group       QueryFieldType = constant.Group
)

var ArrayValues = map[string]QueryFieldType{
	"in":          Array,
	"not_in":      Array,
	"or_in_eq":    OrInEq,
	"or_in_like":  OrInLike,
	"or_out_eq":   OrOutEq,
	"or_out_like": OrOutLike,
}

func (baseHandler *BaseHandler) FieldExtractor(c *gin.Context, model interface{}) {
	var (
		queries []*QueryField
		join    bool
	)

	fields := reflection.TransferFields(model, false)
	fields = AppendJoinFields(c, model, fields)
	for _, field := range fields {
		value, ok := c.GetQueryMap(field)
		if !ok {
			continue
		}
		query := &QueryField{Name: field, Value: value}
		queries = Transform(c, query, queries)
		for _, queryField := range queries {
			if strings.Contains(queryField.Name, constant.Delimiter) {
				JoinExtractor(queryField)
				join = true
			}
		}
	}
	baseHandler.Fields = queries
	baseHandler.Join = join
}

func Transform(c *gin.Context, queryField *QueryField, queries []*QueryField) []*QueryField {
	switch queryField.Value.(type) {
	case map[string]string:
		items, _ := queryField.Value.(map[string]string)
		for k, v := range items {
			field := &QueryField{}
			field.Name = queryField.Name
			field.Operation = k
			field.Type = String
			field.Value = v
			if stringutil.StrInsenContains(k, constant.NotNull, constant.SnotNull, constant.IsNull, constant.IsEmpty) {
				field.Type = NullOrExist
			} else if stringutil.StrInsensitive(constant.Order, k) {
				field.Type = Order
			} else if stringutil.StrInsensitive(constant.Group, k) {
				field.Type = Group
			} else if typ, ok := ArrayValues[k]; ok {
				field.Type = typ
				ParseArrayParams(c, field)
			}
			queries = append(queries, field)
		}
	}
	return queries
}

func ParseArrayParams(c *gin.Context, queryField *QueryField) {
	key := fmt.Sprintf("%s[%s]", queryField.Name, queryField.Operation)
	query := c.Request.URL.Query()
	queryField.Value = query[key]
}

func AppendJoinFields(c *gin.Context, model interface{}, fields []string) []string {
	flag := false
	query := c.Request.URL.Query()
	for key, _ := range query {
		if strings.Contains(key, constant.Delimiter) {
			flag = true
			break
		}
	}
	if flag {
		if method, ok := model.(Joinser); ok {
			for _, join := range method.Joins() {
				fields = append(fields, reflection.TransferFields(join, true)...)
			}
		}
	}
	return fields
}

func JoinExtractor(queryField *QueryField) {
	arr := strings.Split(queryField.Name, constant.Delimiter)
	if len(arr) != 2 {
		panic(constant.QueryParamError)
	}
	joinTable := strcase.ToCamel(arr[0])
	queryField.Join = true
	queryField.JoinTable = joinTable
	queryField.Name = fmt.Sprintf(`"%s"."%s"`, joinTable, arr[1])
}
