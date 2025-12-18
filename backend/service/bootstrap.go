package service

import (
	"context"

	"rungrid/backend/domain"
)

func EnsureDefaultGroups(ctx context.Context, groups *GroupService) error {
	if groups == nil {
		return nil
	}

	existing, err := groups.List(ctx)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	_, err = groups.Create(ctx, domain.GroupInput{
		Name:  "dev",
		Order: 0,
		Color: "#4f7dff",
	})
	return err
}
