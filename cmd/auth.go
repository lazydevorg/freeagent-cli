package cmd

import (
	"github.com/lazydevorg/freeagent-cli/internal/client/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with FreeAgent",
	Run: func(cmd *cobra.Command, args []string) {
		auth.NewOAuthServer().Authenticate()
	},
}
