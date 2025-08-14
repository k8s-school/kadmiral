package cmd

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/example/kadmiral/pkg/remote"
	"github.com/spf13/cobra"
)

var prereqOS string

var prereqCmd = &cobra.Command{
	Use:   "prereq",
	Short: "run prerequisites on nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		hosts := nodeList()
		if len(hosts) == 0 {
			return fmt.Errorf("no nodes specified")
		}
		script := filepath.Join("/tmp/kadmiral/resource", prereqOS, "prereq.sh")
		slog.Info("running prerequisites", "os", prereqOS, "nodes", hosts)
		if err := remote.RunScript(hosts, SSHUser, SSHKey, script); err != nil {
			return err
		}
		slog.Info("prerequisites complete")
		return nil
	},
}

func init() {
	prereqCmd.Flags().StringVar(&prereqOS, "os", "ubuntu", "target operating system")
	rootCmd.AddCommand(prereqCmd)
}
