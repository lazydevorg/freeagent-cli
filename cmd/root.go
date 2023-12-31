package cmd

import (
	"fmt"
	"github.com/lazydevorg/freeagent-cli/cmd/auth"
	"github.com/lazydevorg/freeagent-cli/cmd/whoami"
	"github.com/lazydevorg/freeagent-cli/internal/client"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(auth.Cmd)
	RootCmd.AddCommand(whoami.Cmd)
}

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

func run(_ *cobra.Command, args []string) {
	fmt.Printf("root called: %s\n", args)
}

func Execute() {
	//RootCmd.SetArgs([]string{"whoami"})
	err := RootCmd.Execute()
	if err != nil {
		_ = fmt.Errorf("error executing root command: %s", err)
	}
}
