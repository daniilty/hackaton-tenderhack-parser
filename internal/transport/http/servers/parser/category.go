package parser

import (
	"encoding/json"
	"tenderhack-parser/internal/models"
	"tenderhack-parser/internal/services/parser"

	"github.com/valyala/fasthttp"
)

func getErrorMsgJSON(err error) string {
	if err == nil {
		return `{}`
	}

	return `{"error":"` + err.Error() + `"}`
}

func (a *API) getCategories(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Add("content-type", "application/json")

	categories, err := a.categories.GetCategories(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	bb, err := json.Marshal(categories)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	ctx.Write(bb)
}

func (a *API) parseLogCategory(ctx *fasthttp.RequestCtx) {
	log := &models.Log{}
	ctx.Response.Header.Add("content-type", "application/json")

	err := json.Unmarshal(ctx.Request.Body(), log)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	err = a.svc.ProcessLog(ctx, log)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	ctx.WriteString(`{}`)
}

func (a *API) createCategoryGroup(ctx *fasthttp.RequestCtx) {
	cg := &parser.JoinedCategoryGroup{}
	ctx.Response.Header.Add("content-type", "application/json")

	err := json.Unmarshal(ctx.Request.Body(), cg)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	err = a.svc.CreateCategoryGroup(ctx, cg)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	ctx.WriteString(`{}`)
}

func (a *API) getCategoryGroups(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Add("content-type", "application/json")

	categoryGroups, err := a.svc.GetCategoryGroups(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	bb, err := json.Marshal(categoryGroups)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	ctx.Write(bb)
}
