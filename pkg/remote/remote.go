package remote

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"sync"

	"github.com/k8s-school/kadmiral/resources"
)

func RunParallel(hosts []string, user, key, script string, deps []string) ([]string, []error) {
	slog.Info("run script", "script", script, "hosts", hosts)
	var wg sync.WaitGroup
	errCh := make(chan error, len(hosts))
	outCh := make(chan string, len(hosts))
	for _, h := range hosts {
		host := h
		wg.Add(1)
		go func() {
			defer wg.Done()

			for _, file := range append([]string{script}, deps...) {
				if err := uploadScript(host, user, key, file); err != nil {
					slog.Error("failed to upload script", "host", host, "script", file, "err", err)
					errCh <- err
					return
				}
			}

			out, err := RunScript(host, user, key, script, deps)
			msg := strings.TrimSpace(string(out))
			outCh <- fmt.Sprintf("host %s: %s", host, msg)
			if err != nil {
				errCh <- fmt.Errorf("ssh %s: %v: %s", host, err, msg)
			}
		}()
	}
	wg.Wait()
	close(errCh)
	close(outCh)
	var outputs []string
	for msg := range outCh {
		outputs = append(outputs, msg)
	}
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return outputs, errs
}

func RunScript(host string, user, key, script string, deps []string) ([]byte, error) {
	sshArgs := []string{}
	if key != "" {
		sshArgs = append(sshArgs, "-i", key)
	}
	scriptPath := "/tmp/kubeadm/" + script
	sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", user, host), "bash", scriptPath)
	slog.Info("exec script", "host", host, "script", scriptPath)
	cmd := exec.Command("ssh", sshArgs...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		msg := strings.TrimSpace(string(out))
		slog.Error("script failed", "host", host, "err", err, "output", msg)
	} else {
		slog.Info("script complete", "host", host)
	}
	return out, err
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
	return nil
}
