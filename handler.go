package easyquery

import (
	"easyquery/tools/constant"
	"easyquery/tools/stringutil"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	QueryParams
}

type QueryParams struct {
	Fields     []*QueryField
	Join       bool
	Pagination Paginater
}

func (query *QueryParams) GetFields() []*QueryField {
	return query.Fields
}

func (query *QueryParams) GetJoin() bool {
	return query.Join
}

func (query *QueryParams) GetPagination() Paginater {
	return query.Pagination
}

func (baseHandler *BaseHandler) Transform(ctx *gin.Context, model interface{}) QueryParamer {
	baseHandler.FieldExtractor(ctx, model)
	baseHandler.Paginate(ctx)
	return baseHandler
}

func (baseHandler *BaseHandler) Paginate(ctx *gin.Context) {
	size := stringutil.Str2Int(ctx.DefaultQuery(constant.SymbolSize, constant.Size))
	page := stringutil.Str2Int(ctx.DefaultQuery(constant.SymbolPage, constant.Page))
	pagination := !(ctx.Query(constant.SymbolPagable) == "false")
	baseHandler.Pagination = NewPagination(size, page, pagination)
}
