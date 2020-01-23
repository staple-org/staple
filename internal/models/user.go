package models

// User defines a user of the system.
type User struct {
	// Email will be used as username.
	Email string `json:"email"`
	// Password
	Password string `json:"password"`
	// Confirm code -- ignore in json
	ConfirmCode string `json:"-"`
	// Maximum number of staples
	MaxStaples int `json:"max_staples"`
}
