package cmd

import (
	"fmt"
	"log/slog"

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
		slog.Info("initializing control plane", "node", host)
		if err := remote.RunScript([]string{host}, SSHUser, SSHKey, "/tmp/kadmiral/resource/init.sh"); err != nil {
			return err
		}
		slog.Info("control plane initialized", "node", host)
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&initNode, "node", "", "control plane node")
	rootCmd.AddCommand(initCmd)
}
