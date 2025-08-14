package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	SSHUser  string
	SSHKey   string
	Nodes    []string
	logLevel string
)

var rootCmd = &cobra.Command{
	Use:   "kadmiral",
	Short: "kadmiral manages kubernetes clusters via SSH",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	var level slog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	if err := rootCmd.Execute(); err != nil {
		logger.Error("command failed", "err", err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&SSHUser, "user", "root", "SSH user")
	rootCmd.PersistentFlags().StringVar(&SSHKey, "key", "", "Path to SSH private key")
	rootCmd.PersistentFlags().StringSliceVar(&Nodes, "nodes", []string{}, "Comma separated list of nodes")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error)")
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
