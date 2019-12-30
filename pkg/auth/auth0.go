package auth

import (
	"context"
	"log"
	"os"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// Authenticator provides auth0 authentication functionality.
type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

// NewAuthenticator returns a new Authenticator.
func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	issuer := os.Getenv("STAPLE_ISSUER")
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}
	clientID := os.Getenv("STAPLE_CLIENT_ID")
	clientSecret := os.Getenv("STAPLE_CLIENT_SECRET")
	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:9998/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}
