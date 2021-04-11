package router

import (
	"easyquery/examples/pkg/user"
)

func userRoutes() {
	resources := router.Group("/users")

	controller := &user.UserHandler{}
	resources.GET("/", controller.List)
	resources.POST("/", controller.Create)
	resources.PUT("/:id", controller.Update)
	resources.DELETE("/:id", controller.Delete)
	resources.GET("/group", controller.Group)
	resources.GET("/find", controller.Find)
}
