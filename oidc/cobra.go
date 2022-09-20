package oidc

import (
	"context"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func LoginCmd(cfg *oauth2.Config, aud, keyFilePath string, onTokenOrErr func(t *oauth2.Token, err error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login with your Google account.",
		Run: func(cmd *cobra.Command, args []string) {
			var ts oauth2.TokenSource
			if keyFilePath != "" {
				var err error
				ts, err = NewGoogleServiceAccountTokenSource(context.Background(), keyFilePath, aud)
				if err != nil {
					onTokenOrErr(nil, err)
					return
				}
			} else {
				ts = NewTokenSource(context.Background(), cfg, aud)
			}
			onTokenOrErr(ts.Token())
		},
	}

	return cmd
}
