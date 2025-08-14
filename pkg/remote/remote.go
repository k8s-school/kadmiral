package remote

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"sync"
)

func Rsync(hosts []string, user, key, localDir, remoteDir string) error {
	slog.Info("rsync start", "hosts", hosts, "source", localDir, "dest", remoteDir)
	var wg sync.WaitGroup
	errCh := make(chan error, len(hosts))
	for _, h := range hosts {
		host := h
		wg.Add(1)
		go func() {
			defer wg.Done()
			args := []string{"-az", localDir + "/", fmt.Sprintf("%s@%s:%s", user, host, remoteDir)}
			if key != "" {
				args = append([]string{"-e", fmt.Sprintf("ssh -i %s", key)}, args...)
			}
			slog.Info("rsyncing", "host", host)
			cmd := exec.Command("rsync", args...)
			slog.Debug("rsync command", "host", host, "args", args)
			out := []byte{}
			err := error(nil)
			if out, err = cmd.CombinedOutput(); err != nil {
				msg := strings.TrimSpace(string(out))
				slog.Error("rsync failed", "host", host, "err", err, "output", msg)
				errCh <- fmt.Errorf("rsync %s: %v: %s", host, err, msg)
				return
			}
			slog.Info("rsync complete", "host", host)
			slog.Debug("rsync output", "host", host, "output", string(out))
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		return err
	}
	return nil
}

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

func RunCommand(hosts []string, user, key, command string) error {
	slog.Info("run command", "command", command, "hosts", hosts)
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
			sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", user, host), "sudo", "bash", "-c", command)
			slog.Info("exec command", "host", host)
			cmd := exec.Command("ssh", sshArgs...)
			if out, err := cmd.CombinedOutput(); err != nil {
				msg := strings.TrimSpace(string(out))
				slog.Error("command failed", "host", host, "err", err, "output", msg)
				errCh <- fmt.Errorf("ssh %s: %v: %s", host, err, msg)
				return
			}
			slog.Info("command complete", "host", host)
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		return err
	}
	return nil
}
