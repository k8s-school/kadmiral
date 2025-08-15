package remote

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"sync"
)

func RunScript(hosts []string, user, key, scriptPath string) error {
	slog.Info("run script", "script", scriptPath, "hosts", hosts)
	var wg sync.WaitGroup
	errCh := make(chan error, len(hosts))
	for _, h := range hosts {
		host := h
		wg.Add(1)
		go func() {
			defer wg.Done()
			sshArgs := []string{}
			if key != "" {
				sshArgs = append(sshArgs, "-i", key)
			}
			sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", user, host), "sudo", "bash", scriptPath)
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
