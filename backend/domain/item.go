package domain

import "time"

type ItemType string

const (
	ItemTypeApp    ItemType = "app"
	ItemTypeURL    ItemType = "url"
	ItemTypeFolder ItemType = "folder"
	ItemTypeDoc    ItemType = "doc"
	ItemTypeSystem ItemType = "system"
)

type Item struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Type        ItemType   `json:"type"`
	IconPath    string     `json:"icon_path"`
	GroupID     string     `json:"group_id"`
	Tags        []string   `json:"tags"`
	Favorite    bool       `json:"favorite"`
	LaunchCount int64      `json:"launch_count"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	Hidden      bool       `json:"hidden"`
}

type ItemInput struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Type     ItemType `json:"type"`
	IconPath string   `json:"icon_path"`
	GroupID  string   `json:"group_id"`
	Tags     []string `json:"tags"`
	Favorite bool     `json:"favorite"`
	Hidden   bool     `json:"hidden"`
}

type ItemUpdate struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Type     ItemType `json:"type"`
	IconPath string   `json:"icon_path"`
	GroupID  string   `json:"group_id"`
	Tags     []string `json:"tags"`
	Favorite bool     `json:"favorite"`
	Hidden   bool     `json:"hidden"`
}

func (t ItemType) IsValid() bool {
	switch t {
	case ItemTypeApp, ItemTypeURL, ItemTypeFolder, ItemTypeDoc, ItemTypeSystem:
		return true
	default:
		return false
	}
}
