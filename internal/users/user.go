package users

import "github.com/staple-org/staple/internal/service"

// User defines a user of the system.
type User struct {
	Email   string
	ID      string
	Staples []*service.Staple
}
