package cmd

import (
	"fmt"

	"github.com/example/kadmiral/pkg/remote"
	"github.com/spf13/cobra"
)

var initNode string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the control plane",
	RunE: func(cmd *cobra.Command, args []string) error {
		host := initNode
		if host == "" {
			hosts := nodeList()
			if len(hosts) == 0 {
				return fmt.Errorf("no nodes specified")
			}
			host = hosts[0]
		}
		return remote.RunScript([]string{host}, SSHUser, SSHKey, "/tmp/kadmiral/resource/init.sh")
	},
}

func init() {
	initCmd.Flags().StringVar(&initNode, "node", "", "control plane node")
	rootCmd.AddCommand(initCmd)
}
