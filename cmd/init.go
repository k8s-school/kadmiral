package cmd

import (
	"fmt"
	"log/slog"

	"github.com/k8s-school/kadmiral/pkg/remote"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the control plane",
	RunE: func(cmd *cobra.Command, args []string) error {
		host := AppConfig.ControlPlaneNodes[0]
		if host == "" {
			hosts := nodeList()
			if len(hosts) == 0 {
				return fmt.Errorf("no nodes specified")
			}
			host = hosts[0]
		}
		slog.Info("initializing control plane", "node", host)
		if _, err := remote.RunParallel([]string{host}, AppConfig.SSHUser, AppConfig.SSHKey, "init.sh", []string{"kubeadm-config.yaml", "tokens.csv", "wait-for-master.sh"}); err != nil {
			return err[0]
		}
		slog.Info("control plane initialized", "node", host)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
