package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/k8s-school/ciux/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Distrib string   `yaml:"distrib"`
	Master  string   `yaml:"master"`
	Nodes   []string `yaml:"nodes"`
	User    string   `yaml:"user"`
	SCP     string   `yaml:"scp"`
	SSH     string   `yaml:"ssh"`
}

var (
	SSHUser   string
	SSHKey    string
	Nodes     []string
	verbosity int
	cfgFile   string
	AppConfig Config
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
		slog.Error("command failed", "err", err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&SSHUser, "user", "root", "SSH user")
	rootCmd.PersistentFlags().StringVar(&SSHKey, "key", "", "Path to SSH private key")
	rootCmd.PersistentFlags().StringSliceVar(&Nodes, "nodes", []string{}, "Comma separated list of nodes")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "Path to config file")

	rootCmd.PersistentFlags().IntVarP(&verbosity, "verbosity", "v", 0, "Verbosity level (-v0 for minimal, -v2 for maximum)")
	cobra.OnInitialize(initLogger)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return loadConfig()
	}
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

func loadConfig() error {
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &AppConfig); err != nil {
		return err
	}
	if SSHUser == "root" && AppConfig.User != "" {
		SSHUser = AppConfig.User
	}
	if len(Nodes) == 0 && len(AppConfig.Nodes) > 0 {
		Nodes = AppConfig.Nodes
	}
	return nil
}

// setUpLogs set the log output ans the log level
func initLogger() {
	log.Init(verbosity)
}
