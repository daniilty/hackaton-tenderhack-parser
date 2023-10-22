package parser

import (
	"encoding/json"
	"errors"
	"tenderhack-parser/internal/models"
	"time"

	"github.com/valyala/fasthttp"
)

type fromTo struct {
	Interval uint32    `json:"interval"`
	From     time.Time `json:"from"`
	To       time.Time `json:"to"`
}

type idJSON struct {
	ID int `json:"id"`
}

func (a *API) resolveLog(ctx *fasthttp.RequestCtx) {
	req := &idJSON{}
	ctx.Response.Header.Add("content-type", "application/json")

	err := json.Unmarshal(ctx.Request.Body(), req)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	err = a.logs.Resolve(ctx, req.ID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	ctx.WriteString(`{}`)
}

func (a *API) unresolveLog(ctx *fasthttp.RequestCtx) {
	req := &idJSON{}
	ctx.Response.Header.Add("content-type", "application/json")

	err := json.Unmarshal(ctx.Request.Body(), req)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	err = a.logs.Unresolve(ctx, req.ID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	ctx.WriteString(`{}`)
}

func (a *API) getLogs(ctx *fasthttp.RequestCtx) {
	req := &fromTo{}
	ctx.Response.Header.Add("content-type", "application/json")

	resolved := ctx.Request.URI().QueryArgs().GetBool("resolved")
	err := json.Unmarshal(ctx.Request.Body(), req)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	if req.From.After(req.To) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(getErrorMsgJSON(errors.New("from must be <= to")))
		return
	}

	var logs []*models.Log

	if resolved {
		logs, err = a.logs.GetResolved(ctx, req.From, req.To)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.WriteString(getErrorMsgJSON(err))
			return
		}
	} else {
		logs, err = a.logs.GetUnresolved(ctx, req.From, req.To)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.WriteString(getErrorMsgJSON(err))
			return
		}
	}

	bb, err := json.Marshal(logs)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	ctx.Write(bb)
}

func (a *API) getLogsByDate(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Add("content-type", "application/json")

	resolved := ctx.Request.URI().QueryArgs().GetBool("resolved")

	interval, err := ctx.Request.URI().QueryArgs().GetUint("interval")
	if err != nil || interval == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(getErrorMsgJSON(errors.New("you must provide interval")))
		return
	}

	fromBB := ctx.Request.URI().QueryArgs().Peek("from")
	if len(fromBB) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(getErrorMsgJSON(errors.New("you must provide from")))
		return
	}

	from, err := time.Parse(time.RFC3339, string(fromBB))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	toBB := ctx.Request.URI().QueryArgs().Peek("to")
	if len(toBB) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(getErrorMsgJSON(errors.New("you must provide to")))
		return
	}

	to, err := time.Parse(time.RFC3339, string(toBB))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	if from.After(to) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(getErrorMsgJSON(errors.New("from must be <= to")))
		return
	}

	logs, err := a.svc.GetLogsByDate(ctx, from, to, resolved, uint32(interval))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	bb, err := json.Marshal(logs)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(getErrorMsgJSON(err))
		return
	}

	ctx.Write(bb)
}
