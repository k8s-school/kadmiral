package remote

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"sync"

	"github.com/k8s-school/kadmiral/resources"
)

func RunScript(hosts []string, user, key, script string, deps []string) error {
	slog.Info("run script", "script", script, "hosts", hosts)
	var wg sync.WaitGroup
	errCh := make(chan error, len(hosts))
	for _, h := range hosts {
		host := h
		wg.Add(1)
		go func() {
			defer wg.Done()

			scriptPath := "/tmp/kubeadm/" + script
			slog.Debug("copy script", "host", host, "path", scriptPath)
			copyCmd := exec.Command(
				"ssh",
				append(
					func() []string {
						if key != "" {
							return []string{"-i", key}
						}
						return []string{}
					}(),
					fmt.Sprintf("%s@%s", user, host),
					fmt.Sprintf("mkdir -p /tmp/kubeadm && cat > %s", scriptPath),
				)...,
			)
			slog.Debug("copy command", "cmd", copyCmd.String())
			data, err := resources.Fs.ReadFile(script)
			if err != nil {
				slog.Error("failed to read script", "script", script, "err", err)
				errCh <- fmt.Errorf("read script %s: %v", script, err)
				return
			}

			copyCmd.Stdin = strings.NewReader(string(data))
			if out, err := copyCmd.CombinedOutput(); err != nil {
				msg := strings.TrimSpace(string(out))
				slog.Error("failed to copy script", "host", host, "err", err, "output", msg)
				errCh <- fmt.Errorf("copy script to %s: %v: %s", host, err, msg)
				return
			}

			sshArgs := []string{}
			if key != "" {
				sshArgs = append(sshArgs, "-i", key)
			}
			sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", user, host), "bash", scriptPath)
			slog.Info("exec script", "host", host)
			cmd := exec.Command("ssh", sshArgs...)
			if out, err := cmd.CombinedOutput(); err != nil {
				msg := strings.TrimSpace(string(out))
				slog.Error("script failed", "host", host, "err", err, "output", msg)
				errCh <- fmt.Errorf("ssh %s: %v: %s", host, err, msg)
				return
			}
			slog.Info("script complete", "host", host)
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		return err
	}
	return nil
}
