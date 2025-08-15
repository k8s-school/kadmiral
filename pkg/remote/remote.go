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

			for _, file := range append([]string{script}, deps...) {
				if err := uploadScript(host, user, key, file, errCh); err != nil {
					slog.Error("failed to upload script", "host", host, "script", file, "err", err)
					errCh <-err
					return
				}
			}

			sshArgs := []string{}
			if key != "" {
				sshArgs = append(sshArgs, "-i", key)
			}
			scriptPath := "/tmp/kubeadm/" + script
			sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", user, host), "bash", scriptPath)
			slog.Info("exec script", "host", host, "script", scriptPath))
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

func uploadScript(host string, user, key, script string) error {

	slog.Debug("copy script", "host", host, "path", script)
	remotePath := "/tmp/kubeadm/" + script
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
			fmt.Sprintf("mkdir -p /tmp/kubeadm/ubuntu && cat > %s", remotePath),
		)...,
	)
	slog.Debug("copy command", "cmd", copyCmd.String())
	data, err := resources.Fs.ReadFile(script)
	if err != nil {
		slog.Error("failed to read script", "script", script, "err", err)
		return fmt.Errorf("read script %s: %v", script, err)
	}

	copyCmd.Stdin = strings.NewReader(string(data))
	if out, err := copyCmd.CombinedOutput(); err != nil {
		msg := strings.TrimSpace(string(out))
		slog.Error("failed to copy script", "host", host, "err", err, "output", msg)
		return fmt.Errorf("copy script to %s: %v: %s", host, err, msg)
	}
	slog.Debug("script copied", "host", host, "path", remotePath)
	return
}
