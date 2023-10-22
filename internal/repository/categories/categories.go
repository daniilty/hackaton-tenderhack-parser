package categories

import (
	"context"
	"database/sql"
	"errors"
	"tenderhack-parser/internal/models"

	"github.com/jmoiron/sqlx"
)

type Categories struct {
	db *sqlx.DB
}

func NewCategories(db *sqlx.DB) *Categories {
	return &Categories{db}
}

func (c *Categories) InsertCategory(ctx context.Context, reg string) (int, error) {
	const q = `insert into categories(reg) values($1) returning id`

	var id int

	err := c.db.QueryRowContext(ctx, q, reg).Scan(&id)

	return id, err
}

func (c *Categories) DeleteCategory(ctx context.Context, id int) error {
	const q = `delete from categories where id=$1`

	_, err := c.db.ExecContext(ctx, q, id)
	return err
}

func (c *Categories) LinkCategory(ctx context.Context, category *models.Category) error {
	const q = `update categories set group_id=:group_id where id=:id`

	_, err := c.db.NamedExecContext(ctx, q, category)
	return err
}

func (c *Categories) GetCategories(ctx context.Context) ([]*models.Category, error) {
	const q = `select categories.id as id, coalesce(categories.group_id, 0) as group_id, categories.reg as reg, coalesce(category_groups.name, '') as group_name from categories left join category_groups on coalesce(categories.group_id, 0)=category_groups.id`

	categories := []*models.Category{}
	err := c.db.SelectContext(ctx, &categories, q)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *Categories) GetCategoriesByGroup(ctx context.Context, id int) ([]*models.Category, error) {
	const q = `select id, coalesce(group_id, 0) as group_id, reg from categories where group_id=$1`

	categories := []*models.Category{}
	err := c.db.SelectContext(ctx, &categories, q, id)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *Categories) GetCategoryByReg(ctx context.Context, reg string) (*models.Category, error) {
	const q = `select id, reg FROM categories WHERE similarity(reg, $1) > 0.2
		order by similarity(reg, $1) desc limit 1`

	category := &models.Category{}
	err := c.db.QueryRowContext(ctx, q, reg).Scan(&category.ID, &category.Reg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return category, nil
}

func (c *Categories) InsertCategoryGroup(ctx context.Context, name string) (int, error) {
	const q = `insert into category_groups(name, severity) values($1, 0) returning id`

	var id int

	err := c.db.QueryRowContext(ctx, q, name).Scan(&id)

	return id, err
}

func (c *Categories) GetCategoryGroups(ctx context.Context) ([]*models.CategoryGroup, error) {
	const q = `select id, name, severity from category_groups`

	cgs := []*models.CategoryGroup{}
	err := c.db.SelectContext(ctx, &cgs, q)
	if err != nil {
		return nil, err
	}

	return cgs, nil
}
