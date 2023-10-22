package parser

import (
	"context"
	"fmt"
	"tenderhack-parser/internal/services/parser"

	"github.com/valyala/fasthttp"
)

type API struct {
	addr       string
	svc        *parser.Service
	server     *fasthttp.Server
	logs       parser.LogRepo
	categories parser.CategoriesRepo
}

func NewAPI(svc *parser.Service, addr string, categories parser.CategoriesRepo, logs parser.LogRepo) *API {
	a := &API{
		svc:        svc,
		addr:       addr,
		logs:       logs,
		categories: categories,
	}

	a.server = &fasthttp.Server{
		Handler: a.router,
	}

	return a
}

func (a *API) router(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Add("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	ctx.Response.Header.Add("Access-Control-Allow-Headers", "*")

	if string(ctx.Method()) == fasthttp.MethodOptions {
		ctx.WriteString(`{}`)
		return
	}

	switch string(ctx.Path()) {
	case "/api/category":
		switch string(ctx.Method()) {
		case fasthttp.MethodPost:
			a.parseLogCategory(ctx)
			return
		}
	case "/api/logs":
		switch string(ctx.Method()) {
		case fasthttp.MethodPost:
			a.parseLogCategory(ctx)
			return
		case fasthttp.MethodGet:
			a.getLogs(ctx)
			return
		}
	case "/api/logs_by_date":
		switch string(ctx.Method()) {
		case fasthttp.MethodGet:
			a.getLogsByDate(ctx)
			return
		}
	case "/api/logs/resolve":
		switch string(ctx.Method()) {
		case fasthttp.MethodPost:
			a.resolveLog(ctx)
			return
		}
	case "/api/logs/unresolve":
		switch string(ctx.Method()) {
		case fasthttp.MethodPost:
			a.unresolveLog(ctx)
			return
		}
	case "/api/categories":
		switch string(ctx.Method()) {
		case fasthttp.MethodGet:
			a.getCategories(ctx)
			return
		}
	case "/api/category_groups":
		switch string(ctx.Method()) {
		case fasthttp.MethodPost:
			a.createCategoryGroup(ctx)
			return
		case fasthttp.MethodGet:
			a.getCategoryGroups(ctx)
			return
		}
	}
	ctx.NotFound()
}

func (a *API) Run(ctx context.Context) {
	fmt.Println("starting api server", "addr", a.addr)
	go func() {
		err := a.server.ListenAndServe(a.addr)
		if err != nil {
			fmt.Println("listen and serve", err.Error())
		}
	}()

	<-ctx.Done()
	fmt.Println("stopping api server")
	err := a.server.Shutdown()
	if err != nil {
		fmt.Println("shutdown server", err.Error())
	}
}
