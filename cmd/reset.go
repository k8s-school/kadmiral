package cmd

import (
	"fmt"

	"github.com/example/kadmiral/pkg/remote"
	"github.com/spf13/cobra"
)

var (
	resetAll bool
)

var resetCmd = &cobra.Command{
	Use:   "reset [node]",
	Short: "reset nodes using reset.sh",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hosts := []string{}
		if resetAll {
			hosts = nodeList()
		} else {
			if len(args) != 1 {
				return fmt.Errorf("node name required unless --all is set")
			}
			hosts = []string{args[0]}
		}
		return remote.RunScript(hosts, SSHUser, SSHKey, "/tmp/kadmiral/resource/reset.sh")
	},
}

func init() {
	resetCmd.Flags().BoolVar(&resetAll, "all", false, "reset all nodes")
	rootCmd.AddCommand(resetCmd)
}
