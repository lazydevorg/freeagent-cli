package timetracking

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "timetracking",
	Short: "Manage time tracking",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
