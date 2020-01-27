package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/staple-org/staple/internal/models"
	"github.com/staple-org/staple/internal/service"
	"github.com/staple-org/staple/internal/storage"
	"github.com/staple-org/staple/pkg/config"
)

func TestListStaples(t *testing.T) {
	inMemoryStapleStore := storage.NewInMemoryStapleStorer()
	stapleHandler := service.NewStapler(inMemoryStapleStore)
	e := echo.New()
	testUser := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "",
		MaxStaples:  25,
	}
	stapleHandler.Create(models.Staple{
		Name:      "TestStaple",
		ID:        0,
		Content:   "TestContent",
		CreatedAt: time.Date(1981, 3, 28, 0, 0, 0, 0, time.UTC),
		Archived:  false,
	}, &testUser)

	config.Opts.GlobalTokenKey = "test"
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = testUser.Email // from context
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	tok, err := token.SignedString([]byte(config.Opts.GlobalTokenKey))
	if err != nil {
		t.Fatal(err)
	}
	t.Run("successful staple list", func(tt *testing.T) {
		req := httptest.NewRequest(echo.GET, "/rest/api/1/staple", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		lister := ListStaples(stapleHandler)
		_ = lister(c)

		if rec.Code != http.StatusOK {
			tt.Fatal("test failed with invalid code: ", rec.Code)
		}
		// Check content
		body, _ := ioutil.ReadAll(rec.Body)
		var staple struct {
			Staples []models.Staple `json:"staples"`
		}
		err = json.Unmarshal(body, &staple)
		assert.Len(tt, staple.Staples, 1, "should have returned a single result")
		assert.Equal(tt, "TestStaple", staple.Staples[0].Name, "expected body did not match")
		assert.Equal(tt, "TestContent", staple.Staples[0].Content, "expected body did not match")
	})
}

func TestAddStaples(t *testing.T) {
	inMemoryStapleStore := storage.NewInMemoryStapleStorer()
	stapleHandler := service.NewStapler(inMemoryStapleStore)
	inMemoryUserStore := storage.NewInMemoryUserStorer()
	notifier := service.NewBufferNotifier()
	userHandler := service.NewUserHandler(context.Background(), inMemoryUserStore, notifier)

	e := echo.New()
	testUser := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "",
		MaxStaples:  25,
	}
	userHandler.Register(testUser)

	config.Opts.GlobalTokenKey = "test"
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = testUser.Email // from context
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	tok, err := token.SignedString([]byte(config.Opts.GlobalTokenKey))
	if err != nil {
		t.Fatal(err)
	}
	t.Run("successful staple add", func(tt *testing.T) {
		req := httptest.NewRequest(echo.POST, "/rest/api/1/staple", bytes.NewBuffer([]byte(`{"name": "testcreate", "content":"testcontent"}`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		adder := AddStaple(stapleHandler, userHandler)
		err = adder(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)

		// Check if staple was created
		req = httptest.NewRequest(echo.GET, "/rest/api/1/staple/0", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("0")
		getter := GetStaple(stapleHandler)
		err = getter(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)

		// Check content
		body, _ := ioutil.ReadAll(rec.Body)
		var staple struct {
			Staple models.Staple `json:"staple"`
		}
		err = json.Unmarshal(body, &staple)
		assert.NoError(tt, err)
		log.Println("The body: ", string(body))
		assert.Equal(tt, "testcreate", staple.Staple.Name)
		assert.Equal(tt, "testcontent", staple.Staple.Content)
	})
}

func TestDeleteStaples(t *testing.T) {
	inMemoryStapleStore := storage.NewInMemoryStapleStorer()
	stapleHandler := service.NewStapler(inMemoryStapleStore)
	e := echo.New()
	testUser := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "",
		MaxStaples:  25,
	}
	stapleHandler.Create(models.Staple{
		Name:      "TestStaple",
		ID:        0,
		Content:   "TestContent",
		CreatedAt: time.Date(1981, 3, 28, 0, 0, 0, 0, time.UTC),
		Archived:  false,
	}, &testUser)

	config.Opts.GlobalTokenKey = "test"
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = testUser.Email // from context
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	tok, err := token.SignedString([]byte(config.Opts.GlobalTokenKey))
	if err != nil {
		t.Fatal(err)
	}
	t.Run("successful staple delete", func(tt *testing.T) {
		// Check if staple was created
		req := httptest.NewRequest(echo.DELETE, "/rest/api/1/staple/0", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("0")
		deleter := DeleteStaple(stapleHandler)
		err = deleter(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)

		// Check that it is really gone
		req = httptest.NewRequest(echo.GET, "/rest/api/1/staple/0", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("0")
		getter := GetStaple(stapleHandler)
		err = getter(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestArchiveStaples(t *testing.T) {
	inMemoryStapleStore := storage.NewInMemoryStapleStorer()
	stapleHandler := service.NewStapler(inMemoryStapleStore)
	e := echo.New()
	testUser := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "",
		MaxStaples:  25,
	}
	stapleHandler.Create(models.Staple{
		Name:      "TestStaple",
		ID:        0,
		Content:   "TestContent",
		CreatedAt: time.Date(1981, 3, 28, 0, 0, 0, 0, time.UTC),
		Archived:  false,
	}, &testUser)

	config.Opts.GlobalTokenKey = "test"
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = testUser.Email // from context
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	tok, err := token.SignedString([]byte(config.Opts.GlobalTokenKey))
	if err != nil {
		t.Fatal(err)
	}
	t.Run("successful staple archive", func(tt *testing.T) {
		// Check if staple was created
		req := httptest.NewRequest(echo.POST, "/rest/api/1/staple/0/archive", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("0")
		archiver := ArchiveStaple(stapleHandler)
		err = archiver(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)

		// Check that it is really gone
		req = httptest.NewRequest(echo.GET, "/rest/api/1/staple/0", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("0")
		getter := GetStaple(stapleHandler)
		err = getter(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)

		// Check that the archive list does contain our staple.
		req = httptest.NewRequest(echo.GET, "/rest/api/1/staple/archive", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		archiveList := ShowArchive(stapleHandler)
		err = archiveList(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		var list struct {
			Staples []models.Staple `json:"staples"`
		}
		body, _ := ioutil.ReadAll(rec.Body)
		err = json.Unmarshal(body, &list)
		if err != nil {
			tt.Fatal(err)
		}
		assert.Len(tt, list.Staples, 1, "should have exactly one element")
		assert.Equal(tt, "TestContent", list.Staples[0].Content)
	})
}
