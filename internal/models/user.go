package models

// User defines a user of the system.
type User struct {
	// Auth0 provides unique nicknames. It's safe to use them as an ID.
	Username string `json:"username"`
}
