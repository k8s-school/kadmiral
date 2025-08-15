package cmd

import (
	"fmt"
	"log/slog"

	"github.com/k8s-school/kadmiral/pkg/remote"
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
		slog.Info("upgrading cluster", "nodes", hosts)
		if err := remote.RunScript(hosts, SSHUser, SSHKey, "upgrade.sh", nil); err != nil {
			return err
		}
		slog.Info("upgrade complete", "nodes", hosts)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
