package auth

import (
	"github.com/lazydevorg/freeagent-cli/internal/client/auth"
	"github.com/spf13/cobra"
)

var force bool

var Cmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with FreeAgent",
	Run: func(cmd *cobra.Command, args []string) {
		auth.Authenticate(force)
	},
}

func init() {
	Cmd.Flags().BoolVar(&force, "force", false, "Force re-authentication")
}
