package cmd

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/k8s-school/kadmiral/pkg/remote"
	"github.com/spf13/cobra"
)

var prereqOS string

var prereqCmd = &cobra.Command{
	Use:   "prereq",
	Short: "run prerequisites on nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		hosts := append(AppConfig.ControlPlaneNodes, AppConfig.WorkerNodes...)
		if len(hosts) == 0 {
			return fmt.Errorf("no nodes specified")
		}
		slog.Info("running prerequisites", "os", prereqOS, "nodes", hosts)
		script := filepath.Join(prereqOS, "prereq.sh")
		if _, err := remote.RunParallel(hosts, AppConfig.SSHUser, AppConfig.SSHKey, script, []string{"env.sh"}); err != nil {
			return fmt.Errorf("failed to run prerequisites: %v", err)
		}
		slog.Info("prerequisites complete")
		return nil
	},
}

func init() {
	prereqCmd.Flags().StringVar(&prereqOS, "os", "ubuntu", "target operating system")
	rootCmd.AddCommand(prereqCmd)
}
