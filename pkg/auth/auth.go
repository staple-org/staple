package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// Middleware defines an authentication middleware using Auth 0.
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("auth-session", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "User not logged in.")
		}
		if _, ok := sess.Values["profile"]; !ok {
			c.Error(err)
			return echo.NewHTTPError(http.StatusUnauthorized, "User not logged in.")
		}
		return next(c)
	}
}

// LoginHandler creates a session with a unique value and redirects to the callback handler.
func LoginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Generate random state
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			c.Error(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		state := base64.StdEncoding.EncodeToString(b)

		sess, err := session.Get("auth-session", c)
		if err != nil {
			c.Error(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		sess.Values["state"] = state
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			c.Error(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		authenticator, err := NewAuthenticator()
		if err != nil {
			c.Error(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.Redirect(http.StatusTemporaryRedirect, authenticator.Config.AuthCodeURL(state))
	}
}

// Callback handles the callback from Auth0 service.
func Callback() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("auth-session", c)
		if err != nil {
			c.Error(err)
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		if c.QueryParam("state") != sess.Values["state"] {
			err := errors.New("state did not match with stored session value")
			c.Error(err)
			log.Println(err)
			//log.Println("Query: ", c.QueryParam("state"))
			//log.Println("Session: ", sess.Values["state"])
			return c.NoContent(http.StatusBadRequest)
		}

		authenticator, err := NewAuthenticator()
		if err != nil {
			c.Error(err)
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		token, err := authenticator.Config.Exchange(context.Background(), c.QueryParam("code"))
		if err != nil {
			log.Println("No token found.")
			c.Error(err)
			log.Println(err)
			return c.NoContent(http.StatusUnauthorized)
		}

		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			err := errors.New("no id_token field in oauth2 token")
			c.Error(err)
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		clientID := os.Getenv("STAPLE_CLIENT_ID")
		oidcConfig := &oidc.Config{
			ClientID: clientID,
		}

		idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)
		if err != nil {
			c.Error(err)
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			c.Error(err)
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		sess.Values["id_token"] = rawIDToken
		sess.Values["access_token"] = token.AccessToken
		sess.Values["profile"] = profile
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			c.Error(err)
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		// Hande auth zero callback.
		// TODO: This won't be correct... Make a nice landing page of some sort?
		url := c.Scheme() + "://" + c.Request().Host + "/rest/api/1/staple/"
		return c.Redirect(http.StatusTemporaryRedirect, url)
	}
}
