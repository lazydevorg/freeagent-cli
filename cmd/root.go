package cmd

import (
	"fmt"
	"github.com/lazydevorg/freeagent-cli/cmd/auth"
	"github.com/lazydevorg/freeagent-cli/cmd/time"
	"github.com/lazydevorg/freeagent-cli/cmd/whoami"
	"github.com/lazydevorg/freeagent-cli/internal/client"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "freeagent",
	Short: "FreeAgent CLI",
	Long:  `FreeAgent CLI Long description`,
	Args:  cobra.OnlyValidArgs,
	Run:   run,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		client.SaveToken()
	},
}

func init() {
	RootCmd.AddCommand(auth.Cmd)
	RootCmd.AddCommand(whoami.Cmd)
	RootCmd.AddCommand(time.Cmd)
}

func run(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}

func Execute() {
	RootCmd.SetArgs([]string{"time", "week"})
	err := RootCmd.Execute()
	if err != nil {
		_ = fmt.Errorf("error executing root command: %s", err)
	}
}
