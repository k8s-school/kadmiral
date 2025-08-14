package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	SSHUser string
	SSHKey  string
	Nodes   []string
)

var rootCmd = &cobra.Command{
	Use:   "kadmiral",
	Short: "kadmiral manages kubernetes clusters via SSH",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&SSHUser, "user", "root", "SSH user")
	rootCmd.PersistentFlags().StringVar(&SSHKey, "key", "", "Path to SSH private key")
	rootCmd.PersistentFlags().StringSliceVar(&Nodes, "nodes", []string{}, "Comma separated list of nodes")
}

func nodeList() []string {
	var list []string
	for _, n := range Nodes {
		if trimmed := strings.TrimSpace(n); trimmed != "" {
			list = append(list, trimmed)
		}
	}
	return list
}
