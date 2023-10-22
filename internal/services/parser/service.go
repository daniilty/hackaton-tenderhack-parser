package parser

import (
	"context"
	"strings"
	"tenderhack-parser/internal/models"
)

const databaseQueryError = "db_query_error_shitfuck"

var dbKeywords = []string{
	"database",
	"sql",
	"db",
	"query",
	"table",
	"constraint",
	"foreign key",
	"deadlock",
}

type Service struct {
	categories CategoriesRepo
	logs       LogRepo
}

func NewService(categories CategoriesRepo, logs LogRepo) *Service {
	return &Service{
		categories: categories,
		logs:       logs,
	}
}

func (s *Service) ProcessLog(ctx context.Context, log *models.Log) error {
	stripped := stripLog(log.Data)

	category, err := s.categories.GetCategoryByReg(ctx, stripped)
	if err != nil {
		return err
	}

	var categoryID int
	if category == nil {
		categoryID, err = s.categories.InsertCategory(ctx, stripped)
		if err != nil {
			return err
		}
	} else {
		categoryID = category.ID
	}

	log.CategoryID = categoryID

	err = s.logs.UpsertLog(ctx, log)
	if err != nil {
		return err
	}

	return nil
}

func stripLog(log string) string {
	lowerLog := strings.ToLower(log)
	for _, kw := range dbKeywords {
		if strings.Contains(lowerLog, kw) {
			return databaseQueryError
		}
	}

	if pos := strings.IndexByte(log, ':'); pos != -1 {
		return log[:pos]
	}

	if pos := strings.IndexByte(log, '.'); pos != -1 {
		return log[:pos]
	}

	return log
}
