package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/k8s-school/kadmiral/pkg/remote"
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
		env := fmt.Sprintf("DISTRIB=\"%s\"\nMASTER=\"%s\"\nNODES=\"%s\"\nUSER=%s\nSCP=\"%s\"\nSSH=\"%s\"\n",
			AppConfig.Distrib, AppConfig.Master, strings.Join(AppConfig.Nodes, " "), AppConfig.User, AppConfig.SCP, AppConfig.SSH)
		if err := os.WriteFile("env.sh", []byte(env), 0644); err != nil {
			return err
		}
		defer os.Remove("env.sh")
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
