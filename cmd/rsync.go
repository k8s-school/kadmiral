package cmd

import (
	"fmt"
	"log/slog"

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
		slog.Info("uploading repository", "nodes", hosts)
		if err := remote.Rsync(hosts, SSHUser, SSHKey, ".", "/tmp/kadmiral"); err != nil {
			return err
		}
		slog.Info("upload complete")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rsyncCmd)
}
