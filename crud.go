package easyquery

import (
	"easyquery/tools/constant"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

type Crud interface {
	List(models interface{}, queries QueryParamer) error
	Find(model interface{}, queries QueryParamer) error
	PrepareScope(queries QueryParamer) *gorm.DB
	Group(models interface{}, queries QueryParamer) error
	Create(model interface{}) error
	Update(model interface{}) error
	Delete(model interface{}) error
}

type CrudService struct {
	Model  interface{}
	Models interface{}
	Table  string
	GormDB func() *gorm.DB
}

func NewCrudService(model, models interface{}, gormDB func() *gorm.DB) *CrudService {
	return &CrudService{model,
		models,
		ModelTableName(model),
		gormDB,
	}
}

func (service *CrudService) List(models interface{}, queries QueryParamer) error {
	var count int64
	pagination := queries.GetPagination()
	currentScope := service.Search(queries).Model(service.Model)
	currentScope.Count(&count)
	pagination.SetTotal(count)
	currentScope = service.Preload(currentScope)
	return service.Paginate(currentScope, models, pagination).Error
}

func (service *CrudService) Paginate(currentScope *gorm.DB, models interface{}, paginater Paginater) *gorm.DB {
	if paginater.GetPagable() {
		currentScope = currentScope.Scopes(PageOrderIdDescScope(paginater, service.Table))
	}
	return currentScope.Find(models)
}

func (service *CrudService) PrepareScope(queries QueryParamer) *gorm.DB {
	return service.Preload(service.Search(queries).Model(service.Model))
}

func (service *CrudService) Group(models interface{}, queries QueryParamer) error {
	currentScope := service.Search(queries).Model(service.Model)
	currentScope.Config.DryRun = true
	currentScope.Find(service.Models)
	currentScope.DryRun = false
	// 避免join预加载查询select默认拼接join对象字段报错的问题，Company__id..
	currentScope.Statement.SQL.Reset()
	currentScope.Statement.Vars = nil
	groupField := funk.Find(queries.GetFields(), func(field *QueryField) bool {
		return field.Type == Group
	})
	if queryField, ok := groupField.(*QueryField); ok {
		return currentScope.Scopes(GroupScope(queryField.Name)).Find(models).Error
	}
	return errors.New(constant.GroupParamError)
}

func (service *CrudService) Create(model interface{}) error {
	return service.GormDB().Model(service.Model).Create(model).Error
}

func (service *CrudService) Update(model interface{}) error {
	return service.GormDB().Model(model).Updates(model).Error
}

func (service *CrudService) Delete(model interface{}) error {
	return service.GormDB().Model(model).Delete(model).Error
}

func (service *CrudService) Find(model interface{}, queries QueryParamer) error {
	return service.PrepareScope(queries).First(model).Error
}

func (service *CrudService) Preload(currentScope *gorm.DB) *gorm.DB {
	if method, ok := service.Model.(Preloader); ok {
		for _, preload := range method.Preload() {
			currentScope = currentScope.Preload(preload)
		}
	}
	return currentScope
}
