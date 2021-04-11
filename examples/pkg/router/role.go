package router

import (
	"easyquery/examples/pkg/user"
)

func roleRoutes() {
	resources := router.Group("/roles")

	controller := &user.RoleHandler{}
	resources.GET("/", controller.List)
	resources.POST("/", controller.Create)
	resources.PUT("/:id", controller.Update)
	resources.DELETE("/:id", controller.Delete)
}
