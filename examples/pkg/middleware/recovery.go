package middleware

import (
	"easyquery/examples/pkg/errors"

	"github.com/ztrue/tracerr"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			var apiError *errors.APIError
			if err := recover(); err != nil {
				switch err.(type) {
				case *errors.APIError:
					apiError = err.(*errors.APIError)
				case error:
					apiError = errors.CustomError(err.(error).Error())
				default:
					apiError = errors.CustomError(err.(string))
				}
				tracerr.Errorf("ERROR: %s\n", apiError.Msg)
				c.JSON(apiError.Code, apiError)
				return
			}
		}()

		c.Next()
	}
}
