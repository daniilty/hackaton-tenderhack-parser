package parser

import (
	"context"
	"tenderhack-parser/internal/models"
)

type JoinedCategoryGroup struct {
	ID         int                `json:"id"`
	Name       string             `json:"name"`
	Severity   int                `json:"severity"`
	Categories []*models.Category `json:"categories"`
}

func (s *Service) GetCategoryGroups(ctx context.Context) ([]*JoinedCategoryGroup, error) {
	groups, err := s.categories.GetCategoryGroups(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*JoinedCategoryGroup, 0, len(groups))
	for _, g := range groups {
		categories, err := s.categories.GetCategoriesByGroup(ctx, g.ID)
		if err != nil {
			return nil, err
		}

		res = append(res, &JoinedCategoryGroup{
			ID:         g.ID,
			Name:       g.Name,
			Severity:   g.Severity,
			Categories: categories,
		})
	}

	return res, nil
}

func (s *Service) CreateCategoryGroup(ctx context.Context, cg *JoinedCategoryGroup) error {
	if cg == nil {
		return nil
	}

	id, err := s.categories.InsertCategoryGroup(ctx, cg.Name)
	if err != nil {
		return err
	}

	for _, c := range cg.Categories {
		c.GroupID = id
		err = s.categories.LinkCategory(ctx, c)
		if err != nil {
			return err
		}
	}

	return nil
}
