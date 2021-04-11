package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	DatabaseError     = newAPIError(http.StatusInternalServerError, "Database Error")
	ServerError       = newAPIError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	RecordNotFound    = newAPIError(http.StatusNotFound, "Record Not Found")
	MethodNotFound    = newAPIError(http.StatusNotFound, http.StatusText(http.StatusMethodNotAllowed))
	RouteNotFound     = newAPIError(http.StatusNotFound, "Route Not Found")
	UnauthorizedError = newAPIError(http.StatusUnauthorized, "User Unauthorized")
)

type APIError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func newAPIError(code int, msg string) *APIError {
	return &APIError{
		Code: code,
		Msg:  msg,
	}
}

func CustomError(message string) *APIError {
	return newAPIError(http.StatusInternalServerError, message)
}

func HandleNotMethodFound(c *gin.Context) {
	c.JSON(MethodNotFound.Code, MethodNotFound)
}

func HandleNotRouteFound(c *gin.Context) {
	c.JSON(RouteNotFound.Code, RouteNotFound)
}
