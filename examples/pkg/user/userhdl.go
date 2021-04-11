package user

import (
	"easyquery"
	"easyquery/tools/stringutil"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	easyquery.BaseHandler
}

func (handler *UserHandler) List(c *gin.Context) {
	var models []User
	err := UserCrudService.List(&models, handler.Transform(c, Model))
	handler.HandleList(c, &models, err)
}

func (handler *UserHandler) Group(c *gin.Context) {
	var models []easyquery.GroupVO
	err := UserCrudService.Group(&models, handler.Transform(c, Model))
	handler.Handle(c, &models, err)
}

func (handler *UserHandler) Create(c *gin.Context) {
	var model User
	c.ShouldBind(&model)
	err := UserCrudService.Create(&model)
	handler.Handle(c, &model, err)
}

func (handler *UserHandler) Update(c *gin.Context) {
	var model User
	c.ShouldBind(&model)
	model.ID = stringutil.Str2Uint(c.Param("id"))
	err := UserCrudService.Update(&model)
	handler.Handle(c, &model, err)
}

func (handler *UserHandler) Find(c *gin.Context) {
	model := User{}
	err := UserCrudService.Find(&model, handler.Transform(c, Model))
	handler.Handle(c, &model, err)
}

func (handler *UserHandler) Delete(c *gin.Context) {
	model := User{}
	model.ID = stringutil.Str2Uint(c.Param("id"))
	err := UserCrudService.Delete(&model)
	handler.Handle(c, &model, err)
}
