package whoami

import (
	"fmt"
	"github.com/lazydevorg/freeagent-cli/internal/client"
	"github.com/lazydevorg/freeagent-cli/internal/client/auth"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "whoami",
	Short:   "Get the current user",
	Example: "# freeagent whoami\nYou are logged in as John Smith (Director)",
	Args:    cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		auth.Authenticate(false)
	},
	Run: func(cmd *cobra.Command, args []string) {
		profile, err := client.PersonalProfile()
		if err != nil {
			panic(err)
		}
		fmt.Printf("You are logged in as %s %s (%s)", profile.FirstName, profile.LastName, profile.Role)
	},
}
