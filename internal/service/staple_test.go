package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/staple-org/staple/internal/models"
	"github.com/staple-org/staple/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestStapler_Create(t *testing.T) {
	store := storage.NewInMemoryStapleStorer()
	stapler := NewStapler(store)
	u := models.User{Email: "test@test.com", MaxStaples: 10}
	staple := models.Staple{
		Name:      "test-staple",
		ID:        0,
		Content:   "test-content",
		CreatedAt: time.Date(1980, 1, 1, 1, 1, 1, 0, time.UTC),
		Archived:  false,
	}
	err := stapler.Create(staple, &u)
	assert.NoError(t, err)
	got, err := stapler.Get(&u, 0)
	assert.NoError(t, err)
	assert.Equal(t, staple, *got)
}

func TestStapler_Delete(t *testing.T) {
	store := storage.NewInMemoryStapleStorer()
	stapler := NewStapler(store)
	u := models.User{Email: "test@test.com", MaxStaples: 10}
	staple := models.Staple{
		Name:      "test-staple",
		ID:        0,
		Content:   "test-content",
		CreatedAt: time.Date(1980, 1, 1, 1, 1, 1, 0, time.UTC),
		Archived:  false,
	}
	err := stapler.Create(staple, &u)
	assert.NoError(t, err)
	err = stapler.Delete(&u, 0)
	assert.NoError(t, err)
	got, err := stapler.Get(&u, 0)
	assert.NoError(t, err)
	assert.Nil(t, got)
}

func TestStapler_List(t *testing.T) {
	store := storage.NewInMemoryStapleStorer()
	stapler := NewStapler(store)
	u := models.User{Email: "test@test.com", MaxStaples: 10}
	staple := models.Staple{
		Name:      "test-staple",
		ID:        0,
		Content:   "test-content",
		CreatedAt: time.Date(1980, 1, 1, 1, 1, 1, 0, time.UTC),
		Archived:  false,
	}
	err := stapler.Create(staple, &u)
	assert.NoError(t, err)
	s2 := staple
	s2.Name = "test-staple-2"
	err = stapler.Create(s2, &u)
	assert.NoError(t, err)
	list, err := stapler.List(&u)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, staple.Name, list[0].Name)
	assert.Equal(t, s2.Name, list[1].Name)
}

func TestStapler_Archive(t *testing.T) {
	store := storage.NewInMemoryStapleStorer()
	stapler := NewStapler(store)
	u := models.User{Email: "test@test.com", MaxStaples: 10}
	staple := models.Staple{
		Name:      "test-staple",
		ID:        0,
		Content:   "test-content",
		CreatedAt: time.Date(1980, 1, 1, 1, 1, 1, 0, time.UTC),
		Archived:  false,
	}
	err := stapler.Create(staple, &u)
	assert.NoError(t, err)
	err = stapler.Archive(&u, 0)
	assert.NoError(t, err)
	got, err := stapler.Get(&u, 0)
	assert.NoError(t, err)
	assert.Nil(t, got)
	archiveList, err := stapler.ShowArchive(&u)
	assert.NoError(t, err)
	assert.Len(t, archiveList, 1)
	assert.Equal(t, staple.Name, archiveList[0].Name)
}

func TestStapler_GetNext(t *testing.T) {
	store := storage.NewInMemoryStapleStorer()
	stapler := NewStapler(store)
	u := models.User{Email: "test@test.com", MaxStaples: 10}
	staple := models.Staple{
		Name:      "test-staple",
		ID:        0,
		Content:   "test-content",
		CreatedAt: time.Date(1980, 1, 1, 1, 1, 1, 0, time.UTC),
		Archived:  false,
	}
	err := stapler.Create(staple, &u)
	assert.NoError(t, err)
	s2 := staple
	s2.CreatedAt = time.Date(1980, 2, 1, 1, 1, 1, 0, time.UTC)
	s2.Name = "test-staple-2"
	s3 := staple
	s3.CreatedAt = time.Date(1980, 3, 1, 1, 1, 1, 0, time.UTC)
	s3.Name = "test-staple-3"
	err = stapler.Create(s2, &u)
	assert.NoError(t, err)
	err = stapler.Create(s3, &u)
	assert.NoError(t, err)
	got, err := stapler.GetNext(&u)
	assert.NoError(t, err)
	assert.Equal(t, staple.Name, got.Name)
}

func TestStapler_Create_Error_MaxStaples(t *testing.T) {
	store := storage.NewInMemoryStapleStorer()
	stapler := NewStapler(store)
	u := models.User{Email: "test@test.com", MaxStaples: 0}
	staple := models.Staple{
		Name:      "test-staple",
		ID:        0,
		Content:   "test-content",
		CreatedAt: time.Date(1980, 1, 1, 1, 1, 1, 0, time.UTC),
		Archived:  false,
	}
	err := stapler.Create(staple, &u)
	assert.EqualError(t, err, "cannot create more staples than 0; current count is: 0")
}

func TestStapler_Create_Error_FromStorage(t *testing.T) {
	store := storage.NewInMemoryStapleStorer()
	store.Err = fmt.Errorf("unable to store staple")
	stapler := NewStapler(store)
	u := models.User{Email: "test@test.com", MaxStaples: 0}
	staple := models.Staple{
		Name:      "test-staple",
		ID:        0,
		Content:   "test-content",
		CreatedAt: time.Date(1980, 1, 1, 1, 1, 1, 0, time.UTC),
		Archived:  false,
	}
	err := stapler.Create(staple, &u)
	assert.EqualError(t, err, "unable to store staple")
}
