package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"rungrid/backend/domain"
	"rungrid/backend/storage"
)

type ruleConfig struct {
	Version string        `json:"version"`
	Groups  []ruleGroup   `json:"groups"`
	Rules   []ruleMapping `json:"rules"`
}

type ruleGroup struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Order    int    `json:"order"`
	Color    string `json:"color"`
	Icon     string `json:"icon"`
}

type ruleMapping struct {
	GroupID string    `json:"group_id"`
	Match   ruleMatch `json:"match"`
}

type ruleMatch struct {
	TargetName []string `json:"target_name"`
}

func ImportGroupRules(ctx context.Context, data []byte, groups *GroupService, items *ItemService) (domain.RuleImportResult, error) {
	if groups == nil || items == nil {
		return domain.RuleImportResult{}, fmt.Errorf("service unavailable")
	}

	var config ruleConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return domain.RuleImportResult{}, err
	}

	existingGroups, err := groups.List(ctx)
	if err != nil {
		return domain.RuleImportResult{}, err
	}

	groupByKey := map[string]domain.Group{}
	for _, group := range existingGroups {
		groupByKey[normalizeRuleKey(group.ID)] = group
	}

	groupIDMap := map[string]string{}
	groupCategoryMap := map[string]string{}
	for key, group := range groupByKey {
		groupIDMap[key] = group.ID
		groupCategoryMap[group.ID] = group.Category
	}

	result := domain.RuleImportResult{}
	for _, group := range config.Groups {
		key := normalizeRuleKey(group.ID)
		if key == "" {
			return result, fmt.Errorf("group id is required")
		}
		name := strings.TrimSpace(group.Name)
		if name == "" {
			return result, fmt.Errorf("group name is required")
		}

		if existing, ok := groupByKey[key]; ok {
			updated, err := groups.Update(ctx, domain.Group{
				ID:       existing.ID,
				Name:     name,
				Order:    group.Order,
				Color:    group.Color,
				Category: group.Category,
				Icon:     group.Icon,
			})
			if err != nil {
				return result, err
			}
			groupIDMap[key] = updated.ID
			groupCategoryMap[updated.ID] = updated.Category
			result.GroupsUpdated++
			continue
		}

		created, err := groups.Create(ctx, domain.GroupInput{
			Name:     name,
			Order:    group.Order,
			Color:    group.Color,
			Category: group.Category,
			Icon:     group.Icon,
		})
		if err != nil {
			return result, err
		}
		groupIDMap[key] = created.ID
		groupCategoryMap[created.ID] = created.Category
		result.GroupsCreated++
	}

	matchToGroup := map[string]string{}
	for _, rule := range config.Rules {
		groupKey := normalizeRuleKey(rule.GroupID)
		if groupKey == "" {
			continue
		}
		groupID, ok := groupIDMap[groupKey]
		if !ok {
			return result, fmt.Errorf("unknown group id: %s", rule.GroupID)
		}

		for _, value := range rule.Match.TargetName {
			targetName := normalizeTargetName(value)
			if targetName == "" {
				continue
			}
			if _, exists := matchToGroup[targetName]; exists {
				continue
			}
			matchToGroup[targetName] = groupID
		}
	}

	if len(matchToGroup) == 0 {
		return result, nil
	}

	itemsList, err := items.List(ctx, storage.ItemFilter{})
	if err != nil {
		return result, err
	}

	for _, item := range itemsList {
		targetName := normalizeTargetName(item.TargetName)
		if targetName == "" {
			continue
		}
		groupID, ok := matchToGroup[targetName]
		if !ok {
			continue
		}
		if item.GroupID == groupID {
			continue
		}
		if category, ok := groupCategoryMap[groupID]; ok && category != "" {
			if !matchesItemCategory(item.Type, category) {
				continue
			}
		}

		if _, err := items.Update(ctx, domain.ItemUpdate{
			ID:       item.ID,
			GroupID:  groupID,
			Favorite: item.Favorite,
			Hidden:   item.Hidden,
		}); err != nil {
			return result, err
		}
		result.ItemsUpdated++
	}

	return result, nil
}

func normalizeRuleKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func normalizeTargetName(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func matchesItemCategory(itemType domain.ItemType, groupCategory string) bool {
	category := normalizeRuleKey(groupCategory)
	switch itemType {
	case domain.ItemTypeApp:
		return category == "app"
	case domain.ItemTypeSystem:
		return category == "system"
	case domain.ItemTypeURL:
		return category == "url"
	case domain.ItemTypeDoc:
		return category == "doc"
	case domain.ItemTypeFolder:
		return category == "folder"
	default:
		return false
	}
}
