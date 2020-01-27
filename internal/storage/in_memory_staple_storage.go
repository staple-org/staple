package storage

import (
	"errors"
	"sort"

	"github.com/staple-org/staple/internal/models"
)

// InMemoryStapleStorer is a storer which uses a map as a storage backend.
type InMemoryStapleStorer struct {
	// email as key
	stapleStore map[string][]models.Staple
	Err         error // can be set to simulate an error
}

// NewInMemoryStapleStorer creates a new in memory storage medium.
func NewInMemoryStapleStorer() InMemoryStapleStorer {
	return InMemoryStapleStorer{stapleStore: make(map[string][]models.Staple)}
}

// Create will create a staple in the underlying in memory storage medium.
func (p InMemoryStapleStorer) Create(staple models.Staple, email string) error {
	if _, ok := p.stapleStore[email]; !ok {
		p.stapleStore[email] = make([]models.Staple, 0)
	}
	sort.SliceStable(p.stapleStore[email], func(i, j int) bool {
		return p.stapleStore[email][i].ID < p.stapleStore[email][j].ID
	})
	newID := 0
	if len(p.stapleStore[email]) > 0 {
		newID = p.stapleStore[email][len(p.stapleStore[email])-1].ID + 1
	}
	staple.ID = newID
	p.stapleStore[email] = append(p.stapleStore[email], staple)
	return p.Err
}

// Delete removes a staple.
func (p InMemoryStapleStorer) Delete(email string, stapleID int) error {
	staples := p.stapleStore[email]
	deleteAt := -1
	for i, s := range staples {
		if s.ID == stapleID {
			deleteAt = i
			break
		}
	}
	if deleteAt == -1 {
		return errors.New("staple not found")
	}
	staples = append(staples[:deleteAt], staples[deleteAt+1:]...)
	p.stapleStore[email] = staples
	return p.Err
}

// Get retrieves a staple.
func (p InMemoryStapleStorer) Get(email string, stapleID int) (*models.Staple, error) {
	for _, s := range p.stapleStore[email] {
		if s.ID == stapleID && !s.Archived {
			return &s, nil
		}
	}
	return nil, p.Err
}

// Oldest will get the oldest staple that is not archived.
func (p InMemoryStapleStorer) Oldest(email string) (*models.Staple, error) {
	oldest := p.stapleStore[email][0]
	for _, s := range p.stapleStore[email] {
		if s.CreatedAt.Before(oldest.CreatedAt) {
			oldest = s
		}
	}
	return &oldest, p.Err
}

// Archive archives a staple.
func (p InMemoryStapleStorer) Archive(email string, stapleID int) error {
	for i, s := range p.stapleStore[email] {
		if s.ID == stapleID {
			s.Archived = true
			p.stapleStore[email][i] = s
			return p.Err
		}
	}
	return p.Err
}

// List gets all the not archived staples for a user. List will not retrieve the content
// since that can possibly be a large text. We only ever retrieve it when that
// specific staple is Get.
func (p InMemoryStapleStorer) List(email string) ([]models.Staple, error) {
	list := make([]models.Staple, 0)
	for _, s := range p.stapleStore[email] {
		if !s.Archived {
			list = append(list, s)
		}
	}
	return list, p.Err
}

// ShowArchive will return the users archived staples ordered by id.
func (p InMemoryStapleStorer) ShowArchive(email string) ([]models.Staple, error) {
	list := make([]models.Staple, 0)
	for _, s := range p.stapleStore[email] {
		if s.Archived {
			list = append(list, s)
		}
	}
	return list, p.Err
}
