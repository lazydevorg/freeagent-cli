package time

import (
	"github.com/lazydevorg/freeagent-cli/internal/client"
	"github.com/lazydevorg/freeagent-cli/internal/client/timeslip"
	"github.com/spf13/cobra"
)

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "Show current week time tracking",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient(cmd.Context())
		timeslips := timeslip.Timeslips(c)
		week, err := timeslips.Week()
		if err != nil {
			return err
		}
		err = timeslips.PrintTable(week, nil)
		if err != nil {
			return err
		}
		return nil
	},
}
