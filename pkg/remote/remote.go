package remote

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func Rsync(hosts []string, user, key, localDir, remoteDir string) error {
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
			cmd := exec.Command("rsync", args...)
			if out, err := cmd.CombinedOutput(); err != nil {
				errCh <- fmt.Errorf("rsync %s: %v: %s", host, err, strings.TrimSpace(string(out)))
			}
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
			cmd := exec.Command("ssh", sshArgs...)
			if out, err := cmd.CombinedOutput(); err != nil {
				errCh <- fmt.Errorf("ssh %s: %v: %s", host, err, strings.TrimSpace(string(out)))
			}
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
			cmd := exec.Command("ssh", sshArgs...)
			if out, err := cmd.CombinedOutput(); err != nil {
				errCh <- fmt.Errorf("ssh %s: %v: %s", host, err, strings.TrimSpace(string(out)))
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		return err
	}
	return nil
}
