package easyquery

import (
	"fmt"

	"gorm.io/gorm"
)

func PageScope(paginater Paginater) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(paginater.GetSize()).Offset(paginater.GetOffset())
	}
}

func PageOrderIdDescScope(paginater Paginater, table string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf(`"%s"."id" desc`, table)).Limit(paginater.GetSize()).Offset(paginater.GetOffset())
	}
}

func PageOrderIdAscScope(paginater Paginater, table string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf(`"%s"."id" asc`, table)).Limit(paginater.GetSize()).Offset(paginater.GetOffset())
	}
}

func GroupScope(field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fmt.Sprintf("%s as label, count(1) as value", field)).Group(field).Order("value desc")
	}
}

func GroupOrderScope(clause, group, order string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(clause).Group(group).Order(order)
	}
}

func GroupOrderLimitScope(clause, group, order string, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(clause).Group(group).Order(order).Limit(limit)
	}
}
