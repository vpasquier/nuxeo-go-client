package nuxeoclient

// User structure from JSON response
type user struct {
	Username        string   `json:"username"`
	EntityType      string   `json:"entity-type"`
	IsAdministrator bool     `json:"isAdministrator"`
	Groups          []string `json:"groups"`
}