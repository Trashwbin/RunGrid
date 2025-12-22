package domain

type Group struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Order    int    `json:"order"`
	Color    string `json:"color"`
	Category string `json:"category"`
	Icon     string `json:"icon"`
}

// GroupInput describes fields required to create a group.
type GroupInput struct {
	Name     string `json:"name"`
	Order    int    `json:"order"`
	Color    string `json:"color"`
	Category string `json:"category"`
	Icon     string `json:"icon"`
}
