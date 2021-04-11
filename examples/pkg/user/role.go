package user

import (
	"gorm.io/gorm"
)

var SearchRole = &Role{}
var SearchRoles = &[]Role{}

type Role struct {
	gorm.Model
	Name  string
	Slug  string
	Users []User `gorm:"many2many:user_roles;"`
}
