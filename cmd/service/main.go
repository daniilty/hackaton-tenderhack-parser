package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tenderhack-parser/internal/config"
	"tenderhack-parser/internal/repository/categories"
	"tenderhack-parser/internal/repository/logs"
	"tenderhack-parser/internal/services/parser"
	parserserver "tenderhack-parser/internal/transport/http/servers/parser"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	f := initFlags()

	cfg, err := config.Read(f.configPath)
	exitOnErr(err)

	db, err := sqlx.Open("pgx", cfg.PostgresDSN)
	exitOnErr(err)

	logsRep := logs.NewLogs(db)
	categoriesRep := categories.NewCategories(db)

	svc := parser.NewService(categoriesRep, logsRep)

	api := parserserver.NewAPI(svc, cfg.HTTPAddr, categoriesRep, logsRep)
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT)
	api.Run(ctx)
}

func exitOnErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
