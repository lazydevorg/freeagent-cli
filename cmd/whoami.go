package cmd

import (
	"fmt"
	"github.com/lazydevorg/freeagent-cli/internal/client"
	"github.com/spf13/cobra"
)

var whoAmICmd = &cobra.Command{
	Use:   "whoami",
	Short: "Get the current user",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewClient(cmd.Context())
		usersClient := c.NewUsersClient()
		defer c.Close()
		profile, err := usersClient.PersonalProfile()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s %s (%s)", profile.FirstName, profile.LastName, profile.Role)
	},
}
