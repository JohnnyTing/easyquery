package easyquery

import (
	"easyquery/tools/stringutil"
	"fmt"
	"reflect"

	mapset "github.com/deckarep/golang-set"

	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

// common search method
func (base *CrudService) Search(queries QueryParamer) *gorm.DB {
	var (
		field        string
		operation    string
		modelTable   string
		value        interface{}
		currentScope *gorm.DB
		joinable     bool
	)
	modelTable, joinable = base.Table, queries.GetJoin()
	joinTables := mapset.NewSet()
	currentScope = base.GormDB()
	for _, queryField := range queries.GetFields() {
		if queryField.Join {
			joinTables.Add(queryField.JoinTable)
			field = queryField.Name
		} else {
			// 数据库字段为snake_case
			if joinable {
				// avoid join query ambiguous error
				field = fmt.Sprintf(`"%s"."%s"`, modelTable, strcase.ToSnake(queryField.Name))
			} else {
				field = strcase.ToSnake(queryField.Name)
			}
		}
		value = queryField.Value
		operation = strcase.ToCamel(queryField.Operation)
		queryField.Operation = operation
		queryField.Name = field
		switch queryField.Type {
		case String:
			strValue, ok := StringConvert(value)
			if !ok {
				continue
			}
			clause := SqlClause(operation, field)
			value = WrapperValue(strValue, operation)
			currentScope = currentScope.Where(clause, value)
		case Array:
			strArr, ok := ArrConvert(value)
			if !ok {
				continue
			}
			clause := SqlClause(operation, field)
			currentScope = currentScope.Where(clause, strArr)
		case NullOrExist:
			clause := SqlClause(operation, field)
			currentScope = currentScope.Where(clause)
		case Order:
			currentScope = currentScope.Order(fmt.Sprintf("%s %s", queryField.Name, queryField.Value))
		case OrInEq, OrInLike:
			strArr, ok := ArrConvert(value)
			if !ok {
				continue
			}
			newScope := base.BuildOrClause(strArr, queryField, operation, field, OrInLike)
			currentScope = currentScope.Where(newScope)
		case OrOutEq, OrOutLike:
			strArr, ok := ArrConvert(value)
			if !ok {
				continue
			}
			newScope := base.BuildOrClause(strArr, queryField, operation, field, OrOutLike)
			currentScope = currentScope.Or(newScope)
		}
	}
	return BuildJoins(joinTables, currentScope)
}

func BuildJoins(set mapset.Set, currentScope *gorm.DB) *gorm.DB {
	for _, elem := range set.ToSlice() {
		if j, ok := elem.(string); ok {
			currentScope = currentScope.Joins(j)
		}
	}
	return currentScope
}

func ArrConvert(value interface{}) ([]string, bool) {
	strArr, ok := value.([]string)
	if !ok {
		return nil, ok
	}
	if stringutil.IsEmptySilce(strArr) {
		return nil, false
	}
	return strArr, true
}

func StringConvert(value interface{}) (string, bool) {
	strValue, ok := value.(string)
	if !ok {
		return "", ok
	}
	if stringutil.IsEmpty(strValue) {
		return "", false
	}
	return strValue, true
}

func (base *CrudService) BuildOrClause(strArr []string, queryField *QueryField, operation, field string, fieldType QueryFieldType) *gorm.DB {
	newScope := base.GormDB()
	clause := SqlClause(operation, field)
	for i, item := range strArr {
		if queryField.Type == fieldType {
			item = WrapperValue(item, operation)
		}
		if i == 0 {
			newScope = newScope.Where(clause, item)
		} else {
			newScope = newScope.Or(clause, item)
		}
	}
	return newScope
}

func WrapperValue(value, operation string) string {
	switch operation {
	case "Like", "OrInLike", "OrOutLike":
		value = fmt.Sprintf("%%%s%%", value)
	}
	return value
}

func SqlClause(operation string, field string) string {
	if method, ok := Clause[operation]; ok {
		return method(field)
	}
	panic(fmt.Sprintf("%s is not supported", operation))
}

var Clause map[string]func(string) string

func init() {
	query := &Query{}
	ty := reflect.TypeOf(query)
	va := reflect.ValueOf(query)
	num := ty.NumMethod()
	Clause = make(map[string]func(string) string, num)
	for i := 0; i < num; i++ {
		name := ty.Method(i).Name
		if fn, ok := va.MethodByName(name).Interface().(func(string) string); ok {
			Clause[name] = fn
		}
	}
}

type Query struct{}

func (q *Query) Eq(field string) string {
	return fmt.Sprintf("%s = ?", field)
}

func (q *Query) OrInEq(field string) string {
	return q.Eq(field)
}

func (q *Query) OrOutEq(field string) string {
	return q.Eq(field)
}

func (q *Query) Gt(field string) string {
	return fmt.Sprintf("%s > ?", field)
}

func (q *Query) Gteq(field string) string {
	return fmt.Sprintf("%s >= ?", field)
}

func (q *Query) Lt(field string) string {
	return fmt.Sprintf("%s < ?", field)
}

func (q *Query) Lteq(field string) string {
	return fmt.Sprintf("%s <= ?", field)
}

func (q *Query) In(field string) string {
	return fmt.Sprintf("%s IN (?)", field)
}

func (q *Query) NotIn(field string) string {
	return fmt.Sprintf("%s NOT IN (?)", field)
}

func (q *Query) Not(field string) string {
	return fmt.Sprintf("%s <> ?", field)
}

func (q *Query) IsNull(field string) string {
	return fmt.Sprintf("%s is null", field)
}

func (q *Query) IsEmpty(field string) string {
	return fmt.Sprintf("%s is null or trim(%s) = ''", field, field)
}

func (q *Query) Like(field string) string {
	return fmt.Sprintf("%s like ?", field)
}

func (q *Query) OrInLike(field string) string {
	return q.Like(field)
}

func (q *Query) OrOutLike(field string) string {
	return q.Like(field)
}

func (q *Query) DateGt(field string) string {
	return fmt.Sprintf("Date(%s) > ?", field)
}

func (q *Query) DateGteq(field string) string {
	return fmt.Sprintf("Date(%s) >= ?", field)
}

func (q *Query) DateLt(field string) string {
	return fmt.Sprintf("Date(%s) < ?", field)
}

func (q *Query) DateLteq(field string) string {
	return fmt.Sprintf("Date(%s) <= ?", field)
}

func (q *Query) NotNull(field string) string {
	return fmt.Sprintf("%s is not null and trim(%s) != ''", field, field)
}

func (q *Query) SNotNull(field string) string {
	return fmt.Sprintf("%s is not null and trim(%s) != '' and %s != '无' and %s != '不涉及'",
		field, field, field, field)
}
