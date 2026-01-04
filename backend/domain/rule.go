package domain

type RuleImportResult struct {
	GroupsCreated int `json:"groups_created"`
	GroupsUpdated int `json:"groups_updated"`
	ItemsUpdated  int `json:"items_updated"`
}
