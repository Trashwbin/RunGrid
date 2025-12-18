package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"rungrid/backend/domain"
	"rungrid/backend/storage"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) List(ctx context.Context, filter storage.ItemFilter) ([]domain.Item, error) {
	query := `
		SELECT id, name, path, type, icon_path, group_id, tags, favorite, launch_count, last_used_at, hidden
		FROM items
	`
	args := []interface{}{}
	conditions := []string{}

	if filter.GroupID != "" && filter.GroupID != "all" {
		conditions = append(conditions, "group_id = ?")
		args = append(args, filter.GroupID)
	}
	if strings.TrimSpace(filter.Query) != "" {
		conditions = append(conditions, "LOWER(name) LIKE '%' || ? || '%' ")
		args = append(args, strings.ToLower(strings.TrimSpace(filter.Query)))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY favorite DESC, COALESCE(last_used_at, 0) DESC, launch_count DESC, name ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []domain.Item{}
	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ItemRepository) Get(ctx context.Context, id string) (domain.Item, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, path, type, icon_path, group_id, tags, favorite, launch_count, last_used_at, hidden
		FROM items WHERE id = ?
	`, id)

	item, err := scanItem(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Item{}, storage.ErrNotFound
		}
		return domain.Item{}, err
	}

	return item, nil
}

func (r *ItemRepository) GetByPath(ctx context.Context, path string) (domain.Item, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, path, type, icon_path, group_id, tags, favorite, launch_count, last_used_at, hidden
		FROM items WHERE LOWER(path) = LOWER(?)
	`, path)

	item, err := scanItem(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Item{}, storage.ErrNotFound
		}
		return domain.Item{}, err
	}

	return item, nil
}

func (r *ItemRepository) SetIconPath(ctx context.Context, id string, iconPath string) error {
	result, err := r.db.ExecContext(ctx, "UPDATE items SET icon_path = ? WHERE id = ?", iconPath, id)
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

func (r *ItemRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	tags, err := encodeTags(item.Tags)
	if err != nil {
		return domain.Item{}, err
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO items (
			id, name, path, type, icon_path, group_id, tags, favorite, launch_count, last_used_at, hidden
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		item.ID,
		item.Name,
		item.Path,
		string(item.Type),
		item.IconPath,
		item.GroupID,
		tags,
		boolToInt(item.Favorite),
		item.LaunchCount,
		timeToUnix(item.LastUsedAt),
		boolToInt(item.Hidden),
	)
	if err != nil {
		return domain.Item{}, err
	}

	return item, nil
}

func (r *ItemRepository) Update(ctx context.Context, item domain.Item) (domain.Item, error) {
	tags, err := encodeTags(item.Tags)
	if err != nil {
		return domain.Item{}, err
	}

	result, err := r.db.ExecContext(ctx, `
		UPDATE items SET
			name = ?,
			path = ?,
			type = ?,
			icon_path = ?,
			group_id = ?,
			tags = ?,
			favorite = ?,
			launch_count = ?,
			last_used_at = ?,
			hidden = ?
		WHERE id = ?
	`,
		item.Name,
		item.Path,
		string(item.Type),
		item.IconPath,
		item.GroupID,
		tags,
		boolToInt(item.Favorite),
		item.LaunchCount,
		timeToUnix(item.LastUsedAt),
		boolToInt(item.Hidden),
		item.ID,
	)
	if err != nil {
		return domain.Item{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return domain.Item{}, err
	}
	if affected == 0 {
		return domain.Item{}, storage.ErrNotFound
	}

	return item, nil
}

func (r *ItemRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM items WHERE id = ?", id)
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

func (r *ItemRepository) IncrementLaunch(ctx context.Context, id string, usedAt time.Time) (domain.Item, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE items
		SET launch_count = launch_count + 1,
			last_used_at = ?
		WHERE id = ?
	`, usedAt.Unix(), id)
	if err != nil {
		return domain.Item{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return domain.Item{}, err
	}
	if affected == 0 {
		return domain.Item{}, storage.ErrNotFound
	}

	return r.Get(ctx, id)
}

type itemScanner interface {
	Scan(dest ...any) error
}

func scanItem(scanner itemScanner) (domain.Item, error) {
	var (
		item     domain.Item
		typeText string
		tagsText string
		favorite int
		hidden   int
		lastUsed sql.NullInt64
	)

	err := scanner.Scan(
		&item.ID,
		&item.Name,
		&item.Path,
		&typeText,
		&item.IconPath,
		&item.GroupID,
		&tagsText,
		&favorite,
		&item.LaunchCount,
		&lastUsed,
		&hidden,
	)
	if err != nil {
		return domain.Item{}, err
	}

	item.Type = domain.ItemType(typeText)
	item.Favorite = favorite == 1
	item.Hidden = hidden == 1
	if lastUsed.Valid {
		usedAt := time.Unix(lastUsed.Int64, 0)
		item.LastUsedAt = &usedAt
	}

	if tagsText != "" {
		if err := json.Unmarshal([]byte(tagsText), &item.Tags); err != nil {
			return domain.Item{}, err
		}
	}

	return item, nil
}

func encodeTags(tags []string) (string, error) {
	if tags == nil {
		return "[]", nil
	}
	data, err := json.Marshal(tags)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func timeToUnix(value *time.Time) interface{} {
	if value == nil {
		return nil
	}
	return value.Unix()
}
