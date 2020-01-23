package pkg

import (
	"bytes"
	"io/ioutil"
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
		CreatedAt: time.Date(1981, 3, 28, 0, 0, 0, 0, time.Local),
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
	t.Run("successful staple create", func(tt *testing.T) {
		req := httptest.NewRequest(echo.GET, "/rest/api/1/staple", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		lister := ListStaples(stapleHandler)
		_ = lister(c)

		if rec.Code != http.StatusOK {
			t.Fatal("test failed with invalid code: ", rec.Code)
		}
		body, _ := ioutil.ReadAll(rec.Body)
		expected := []byte(`{"staples":[{"name":"TestStaple","id":0,"content":"TestContent","created_at":"1981-03-28T00:00:00+01:00","archived":false}]}
`)
		assert.Equal(t, expected, body, "expected body did not match")
	})
}
