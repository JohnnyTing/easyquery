package router

import (
	"easyquery/examples/config"
	"easyquery/examples/pkg/errors"
	"easyquery/examples/pkg/middleware"
	"fmt"

	"github.com/ztrue/tracerr"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

// Run will start the server
func Register() {
	gin.SetMode(config.GlobalProperties.App.Mode)
	router.NoMethod(errors.HandleNotMethodFound)
	router.NoRoute(errors.HandleNotRouteFound)
	router.Use(middleware.Recovery())
	//router.Use(middleware.AuthMiddlewareFunc())
	registerRoutes()
	port := config.GlobalProperties.App.Port
	err := router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		err = tracerr.Errorf("server run failed")
		panic(err)
	}
}

// getRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func registerRoutes() {
	userRoutes()
	roleRoutes()
}
