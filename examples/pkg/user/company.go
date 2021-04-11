package user

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	Name string
}
