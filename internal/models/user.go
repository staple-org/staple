package models

// User defines a user of the system.
type User struct {
	// Email will be used as username.
	Email string `json:"email"`
	// Password
	Password string `json:"password"`
	// Confirm link -- ignore in json
	ConfirmLink string `json:"-"`
}
