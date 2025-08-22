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
		if _, err := remote.RunParallel(hosts, AppConfig.SSHUser, AppConfig.SSHKey, "upgrade.sh", nil); err != nil {
			return fmt.Errorf("failed to upgrade nodes: %v", err)
		}
		slog.Info("upgrade complete", "nodes", hosts)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
