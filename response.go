package easyquery

import (
	"easyquery/tools/constant"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ResponseList struct {
	Response
	Pagination Paginater `json:"pagination"`
}

func NewResponse(code int, data interface{}) *Response {
	return &Response{Code: code, Data: data}
}

func NewFailResponse(code int, msg string) *Response {
	return &Response{Code: code, Msg: msg}
}

func SuccessResponse(data interface{}) *Response {
	return NewResponse(http.StatusOK, data)
}

func FailResponse(msg string) *Response {
	return NewFailResponse(http.StatusBadRequest, msg)
}

func NewResponseList(code int, data interface{}, pagination Paginater) *ResponseList {
	response := &ResponseList{Pagination: pagination}
	response.Code = code
	response.Data = data
	return response
}

func SuccessResponseList(data interface{}, pagination Paginater) *ResponseList {
	return NewResponseList(http.StatusOK, data, pagination)
}

func FailUnauthorizedResponse() *Response {
	return NewFailResponse(http.StatusUnauthorized, constant.Unauthorized)
}

func (*BaseHandler) ResponseErr(ctx *gin.Context, msg string) {
	ctx.Render(http.StatusBadRequest, NewBaseJsonRender(FailResponse(msg)))
}

func (*BaseHandler) ResponseUnauthorizedErr(ctx *gin.Context) {
	ctx.Render(http.StatusUnauthorized, NewBaseJsonRender(FailUnauthorizedResponse()))
}

func (*BaseHandler) ResponseSuccess(ctx *gin.Context, data interface{}) {
	ctx.Render(http.StatusOK, NewBaseJsonRender(SuccessResponse(data)))
}

func (*BaseHandler) ResponseDownloadSuccess(ctx *gin.Context, name string, bytes []byte) {
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", name))
	ctx.Header("Content-Type", "application/octet-stream")

	if _, err := ctx.Writer.Write(bytes); err != nil {
		panic("download failed ")
	}
}

func (base *BaseHandler) ResponseSuccessList(ctx *gin.Context, data interface{}) {
	ctx.Render(http.StatusOK, NewBaseJsonRender(SuccessResponseList(data, base.GetPagination())))
}

func (base *BaseHandler) Handle(ctx *gin.Context, data interface{}, err error) {
	if err != nil {
		base.ResponseErr(ctx, err.Error())
	} else {
		base.ResponseSuccess(ctx, data)
	}
}

func (base *BaseHandler) HandleList(ctx *gin.Context, data interface{}, err error) {
	if err != nil {
		base.ResponseErr(ctx, err.Error())
	} else {
		base.ResponseSuccessList(ctx, data)
	}
}
