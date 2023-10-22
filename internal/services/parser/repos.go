package parser

import (
	"context"
	"tenderhack-parser/internal/models"
	"time"
)

type CategoriesRepo interface {
	InsertCategory(context.Context, string) (int, error)
	DeleteCategory(context.Context, int) error
	LinkCategory(context.Context, *models.Category) error
	InsertCategoryGroup(context.Context, string) (int, error)
	GetCategoryGroups(context.Context) ([]*models.CategoryGroup, error)
	GetCategories(context.Context) ([]*models.Category, error)
	GetCategoriesByGroup(context.Context, int) ([]*models.Category, error)
	GetCategoryByReg(context.Context, string) (*models.Category, error)
}

type LogRepo interface {
	UpsertLog(context.Context, *models.Log) error
	Resolve(context.Context, int) error
	Unresolve(context.Context, int) error
	GetResolved(context.Context, time.Time, time.Time) ([]*models.Log, error)
	GetUnresolved(context.Context, time.Time, time.Time) ([]*models.Log, error)
}
