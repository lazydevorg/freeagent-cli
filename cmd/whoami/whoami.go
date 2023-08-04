package whoami

import (
	"github.com/lazydevorg/freeagent-cli/internal/cli"
	"github.com/lazydevorg/freeagent-cli/internal/client"
	"github.com/lazydevorg/freeagent-cli/internal/client/auth"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "whoami",
	Short: "Get the current user",
	Args:  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		auth.Authenticate(cmd.Context(), false)
	},
	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewClient(cmd.Context())
		profile, err := c.PersonalProfile()
		if err != nil {
			panic(err)
		}
		cli.RenderEntityTable(profile)
	},
}
