package domain

type ScanResult struct {
	Total    int `json:"total"`
	Inserted int `json:"inserted"`
	Skipped  int `json:"skipped"`
}
