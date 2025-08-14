package cmd

import (
	"fmt"
	"log/slog"

	"github.com/example/kadmiral/pkg/remote"
	"github.com/spf13/cobra"
)

var (
	joinAll    bool
	joinMaster string
	joinToken  string
	joinHash   string
)

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "join nodes to the cluster",
}

var joinNodeCmd = &cobra.Command{
	Use:   "node [name]",
	Short: "join a worker node",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if joinMaster == "" || joinToken == "" || joinHash == "" {
			return fmt.Errorf("master, token and ca-hash must be provided")
		}
		hosts := []string{}
		if joinAll {
			hosts = nodeList()
		} else {
			if len(args) != 1 {
				return fmt.Errorf("node name required unless --all is set")
			}
			hosts = []string{args[0]}
		}
		command := fmt.Sprintf("kubeadm join %s --token %s --discovery-token-ca-cert-hash sha256:%s", joinMaster, joinToken, joinHash)
		slog.Info("joining nodes", "nodes", hosts, "master", joinMaster)
		if err := remote.RunCommand(hosts, SSHUser, SSHKey, command); err != nil {
			return err
		}
		slog.Info("nodes joined", "nodes", hosts)
		return nil
	},
}

func init() {
	joinNodeCmd.Flags().BoolVar(&joinAll, "all", false, "join all nodes")
	joinNodeCmd.Flags().StringVar(&joinMaster, "master", "", "control plane endpoint")
	joinNodeCmd.Flags().StringVar(&joinToken, "token", "", "bootstrap token")
	joinNodeCmd.Flags().StringVar(&joinHash, "ca-hash", "", "discovery token CA cert hash")
	joinCmd.AddCommand(joinNodeCmd)
	rootCmd.AddCommand(joinCmd)
}
