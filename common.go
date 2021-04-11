package easyquery

import (
	"reflect"

	"gorm.io/gorm/schema"
)

var NS = schema.NamingStrategy{
	SingularTable: false,
}

func TableName(str string) string {
	return NS.TableName(str)
}

func ModelTableName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return TableName(t.Name())
}
