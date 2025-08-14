package cmd

import (
	"fmt"
	"github.com/example/kadmiral/pkg/remote"
	"github.com/spf13/cobra"
)

var rsyncCmd = &cobra.Command{
	Use:   "rsync",
	Short: "upload scripts to all nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		hosts := nodeList()
		if len(hosts) == 0 {
			return fmt.Errorf("no nodes specified")
		}
		return remote.Rsync(hosts, SSHUser, SSHKey, ".", "/tmp/kadmiral")
	},
}

func init() {
	rootCmd.AddCommand(rsyncCmd)
}
