package cmd

import (
	"fmt"

	"github.com/example/kadmiral/pkg/remote"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade cluster using upgrade.sh",
	RunE: func(cmd *cobra.Command, args []string) error {
		hosts := nodeList()
		if len(hosts) == 0 {
			return fmt.Errorf("no nodes specified")
		}
		return remote.RunScript(hosts, SSHUser, SSHKey, "/tmp/kadmiral/upgrade.sh")
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
