package domain

type HotkeyBinding struct {
	ID   string `json:"id"`
	Keys string `json:"keys"`
}

type HotkeyIssue struct {
	ID     string `json:"id"`
	Keys   string `json:"keys"`
	Reason string `json:"reason"`
}

type HotkeyApplyResult struct {
	Issues []HotkeyIssue `json:"issues"`
}
