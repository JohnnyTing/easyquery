package user

import (
	"easyquery"
	"easyquery/tools/stringutil"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	easyquery.BaseHandler
}

func (handler *RoleHandler) List(c *gin.Context) {
	var models []Role
	err := RoleCrudService.List(&models, handler.Transform(c, SearchRole))
	handler.HandleList(c, &models, err)
}

func (handler *RoleHandler) Create(c *gin.Context) {
	var model Role
	c.ShouldBind(&model)
	err := RoleCrudService.Create(&model)
	handler.Handle(c, &model, err)
}

func (handler *RoleHandler) Update(c *gin.Context) {
	var model Role
	c.ShouldBind(&model)
	model.ID = stringutil.Str2Uint(c.Param("id"))
	err := RoleCrudService.Update(&model)
	handler.Handle(c, &model, err)
}

func (handler *RoleHandler) Delete(c *gin.Context) {
	var model Role
	c.ShouldBind(&model)
	model.ID = stringutil.Str2Uint(c.Param("id"))
	err := RoleCrudService.Delete(&model)
	handler.Handle(c, &model, err)
}
