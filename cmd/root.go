package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(whoAmICmd)
}

var rootCmd = &cobra.Command{
	Use:   "freeagent",
	Short: "FreeAgent CLI",
	Long:  `FreeAgent CLI Long description`,
	Args:  cobra.OnlyValidArgs,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	fmt.Printf("root called: %s\n", args)
}

func Execute() {
	rootCmd.SetArgs([]string{"whoami"})
	err := rootCmd.Execute()
	if err != nil {
		_ = fmt.Errorf("Error executing root command: %s", err)
	}
}
