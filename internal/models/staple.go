package models

import "time"

// Staple defines a Staple in the system.
type Staple struct {
	Name             string    `json:"name"`
	ID               string    `json:"id"`
	Content          string    `json:"content"`
	CreatedTimestamp time.Time `json:"created_timestamp"`
	Archived         bool      `json:"archived"`
}
