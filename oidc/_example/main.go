package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/raystack/salt/oidc"
)

func main() {
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:5454",
		Scopes:       strings.Split(os.Getenv("OIDC_SCOPES"), ","),
	}
	aud := os.Getenv("OIDC_AUDIENCE")
	keyFile := os.Getenv("GOOGLE_SERVICE_ACCOUNT")

	onTokenOrErr := func(t *oauth2.Token, err error) {
		if err != nil {
			log.Fatalf("oidc login failed: %v", err)
		}

		_ = json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"token_type":    t.TokenType,
			"access_token":  t.AccessToken,
			"expiry":        t.Expiry,
			"refresh_token": t.RefreshToken,
			"id_token":      t.Extra("id_token"),
		})
	}

	_ = oidc.LoginCmd(cfg, aud, keyFile, onTokenOrErr).Execute()
}
