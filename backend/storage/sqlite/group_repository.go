package sqlite

import (
	"context"
	"database/sql"

	"rungrid/backend/domain"
	"rungrid/backend/storage"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) List(ctx context.Context) ([]domain.Group, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, display_order, color, category
		FROM groups
		ORDER BY display_order ASC, name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []domain.Group{}
	for rows.Next() {
		var group domain.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.Order, &group.Color, &group.Category); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *GroupRepository) Get(ctx context.Context, id string) (domain.Group, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, display_order, color, category
		FROM groups WHERE id = ?
	`, id)

	var group domain.Group
	if err := row.Scan(&group.ID, &group.Name, &group.Order, &group.Color, &group.Category); err != nil {
		if err == sql.ErrNoRows {
			return domain.Group{}, storage.ErrNotFound
		}
		return domain.Group{}, err
	}

	return group, nil
}

func (r *GroupRepository) Create(ctx context.Context, group domain.Group) (domain.Group, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO groups (id, name, display_order, color, category)
		VALUES (?, ?, ?, ?, ?)
	`, group.ID, group.Name, group.Order, group.Color, group.Category)
	if err != nil {
		return domain.Group{}, err
	}

	return group, nil
}

func (r *GroupRepository) Update(ctx context.Context, group domain.Group) (domain.Group, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE groups SET name = ?, display_order = ?, color = ?, category = ?
		WHERE id = ?
	`, group.Name, group.Order, group.Color, group.Category, group.ID)
	if err != nil {
		return domain.Group{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return domain.Group{}, err
	}
	if affected == 0 {
		return domain.Group{}, storage.ErrNotFound
	}

	return group, nil
}

func (r *GroupRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM groups WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return storage.ErrNotFound
	}

	return nil
}
