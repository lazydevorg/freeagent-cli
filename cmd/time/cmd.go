package time

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "time",
	Short: "Manage time tracking",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(weekCmd)
}
