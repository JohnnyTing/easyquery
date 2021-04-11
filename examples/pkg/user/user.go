package user

import (
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

var (
	Model  = &User{}
	Models = &[]User{}
)

type User struct {
	gorm.Model
	LoginName *string
	UserName  *string
	Mobile    *int32
	Password  *string
	Gender    *string
	CompanyID uint
	Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Roles     []Role  `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (user *User) Preload() []string {
	return []string{clause.Associations}
}

func (user *User) Joins() []interface{} {
	return []interface{}{Company{}}
}
