package user

import (
	"easyquery"
)

var RoleCrudService *RoleService

type RoleService struct {
	easyquery.Crud
}

func init() {
	RoleCrudService = NewDefaultRoleService()
}

func NewRoleService(crud easyquery.Crud) *RoleService {
	return &RoleService{crud}
}

func NewDefaultRoleService() *RoleService {
	return NewRoleService(easyquery.NewCrudService(SearchRole, SearchRoles, retreiveGormDB))
}
