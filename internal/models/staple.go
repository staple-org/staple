package models

import "time"

// Staple defines a Staple in the system.
type Staple struct {
	Name      string    `json:"name"`
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	Archived  bool      `json:"archived"`
}
