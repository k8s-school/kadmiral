package cmd

import (
	"fmt"
	"log/slog"

	"github.com/k8s-school/kadmiral/pkg/remote"
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

		var hosts []string
		if len(args) == 0 {
			hosts = append(AppConfig.ControlPlaneNodes, AppConfig.WorkerNodes...)
		} else {
			hosts = []string{args[0]}
		}

		slog.Info("resetting nodes", "nodes", hosts)
		if _, err := remote.RunParallel(hosts, AppConfig.SSHUser, AppConfig.SSHKey, "reset.sh", nil); err != nil {
			return fmt.Errorf("failed to reset nodes: %v", err)
		}
		slog.Info("reset complete", "nodes", hosts)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
