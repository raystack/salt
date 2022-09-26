package oidc

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
)

const (
	// Values from OpenID Connect.
	scopeOpenID = "openid"
	audienceKey = "audience"

	// Values used in PKCE implementation.
	// Refer https://www.rfc-editor.org/rfc/rfc7636
	pkceS256               = "S256"
	codeVerifierLen        = 32
	codeChallengeKey       = "code_challenge"
	codeVerifierKey        = "code_verifier"
	codeChallengeMethodKey = "code_challenge_method"
)

func NewTokenSource(ctx context.Context, conf *oauth2.Config, audience string) oauth2.TokenSource {
	conf.Scopes = append(conf.Scopes, scopeOpenID)
	return &authHandlerSource{
		ctx:      ctx,
		config:   conf,
		audience: audience,
	}
}

type authHandlerSource struct {
	ctx      context.Context
	config   *oauth2.Config
	audience string
}

func (source *authHandlerSource) Token() (*oauth2.Token, error) {
	stateBytes, err := randomBytes(10)
	if err != nil {
		return nil, err
	}
	actualState := string(stateBytes)

	codeVerifier, codeChallenge, challengeMethod, err := newPKCEParams()
	if err != nil {
		return nil, err
	}

	// Step 1. Send user to authorization page for obtaining consent.
	url := source.config.AuthCodeURL(actualState,
		oauth2.SetAuthURLParam(audienceKey, source.audience),
		oauth2.SetAuthURLParam(codeChallengeKey, codeChallenge),
		oauth2.SetAuthURLParam(codeChallengeMethodKey, challengeMethod),
	)

	code, receivedState, err := browserAuthzHandler(source.ctx, source.config.RedirectURL, url)
	if err != nil {
		return nil, err
	} else if receivedState != actualState {
		return nil, errors.New("state received in redirection does not match")
	}

	// Step 2. Exchange code-grant for tokens (access_token, refresh_token, id_token).
	tok, err := source.config.Exchange(source.ctx, code,
		oauth2.SetAuthURLParam(audienceKey, source.audience),
		oauth2.SetAuthURLParam(codeVerifierKey, codeVerifier),
	)
	if err != nil {
		return nil, err
	}

	idToken, ok := tok.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id_token not found in token response")
	}
	tok.AccessToken = idToken

	return tok, nil
}

// newPKCEParams generates parameters for 'Proof Key for Code Exchange'.
// Refer https://www.rfc-editor.org/rfc/rfc7636#section-4.2
func newPKCEParams() (verifier, challenge, method string, err error) {
	// generate 'verifier' string.
	verifierBytes, err := randomBytes(codeVerifierLen)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	verifier = encode(verifierBytes)

	// generate S256 challenge.
	h := sha256.New()
	h.Write([]byte(verifier))
	challenge = encode(h.Sum(nil))

	return verifier, challenge, pkceS256, nil
}
