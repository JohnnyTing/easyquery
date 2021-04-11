package middleware

import (
	"easyquery"
	"easyquery/examples/pkg/user"
	"easyquery/tools/stringutil"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
)

var AuthIgnoreList = mapset.NewSet(
	"/login",
	"/auth/token",
)

func AuthMiddlewareFunc() gin.HandlerFunc {
	baseHandler := easyquery.BaseHandler{}
	return func(ctx *gin.Context) {
		var token string
		if AuthIgnoreList.Contains(ctx.Request.RequestURI) {
			ctx.Next()
			return
		}
		token = strings.TrimSpace(ctx.Request.Header.Get("token"))
		if stringutil.IsEmpty(token) {
			token = strings.TrimSpace(ctx.Query("token"))
		}
		if stringutil.IsEmpty(token) {
			ctx.Abort()
			baseHandler.ResponseUnauthorizedErr(ctx)
			return
		}

		find, ok := user.UserCrudService.FindByToken(token)
		if ok != nil {
			ctx.Abort()
			baseHandler.ResponseUnauthorizedErr(ctx)
			return
		}
		ctx.Set("user", find)
	}
}
