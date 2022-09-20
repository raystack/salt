package oidc

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

func NewGoogleServiceAccountTokenSource(ctx context.Context, keyFile, aud string) (oauth2.TokenSource, error) {
	return idtoken.NewTokenSource(ctx, aud, idtoken.WithCredentialsFile(keyFile))
}
