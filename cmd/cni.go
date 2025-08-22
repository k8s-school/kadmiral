package cmd

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/k8s-school/kadmiral/pkg/remote"
	"github.com/spf13/cobra"
)

var cniCmd = &cobra.Command{
	Use:   "cni",
	Short: "manage CNI plugins",
}

var cniInstallCmd = &cobra.Command{
	Use:   "install [name]",
	Short: "install a CNI plugin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		script := filepath.Join("/tmp/kadmiral/resource", fmt.Sprintf("install-%s.sh", name))
		hosts := nodeList()
		if len(hosts) == 0 {
			return fmt.Errorf("no nodes specified")
		}
		slog.Info("installing CNI", "name", name, "node", hosts[0])
		// assume CNI is installed on control plane first node
		if _, err := remote.RunParallel([]string{hosts[0]}, AppConfig.SSHUser, AppConfig.SSHKey, script, nil); err != nil {
			return err[0]
		}
		slog.Info("CNI installed", "name", name)
		return nil
	},
}

func init() {
	cniCmd.AddCommand(cniInstallCmd)
	rootCmd.AddCommand(cniCmd)
}
