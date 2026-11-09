package git

import (
	"argus-agent/internal/config"
	"fmt"
	"os"
	"os/exec"
)

// type repo struct {
// 	name          string
// 	url           string
// 	branches      []string
// 	defaultBranch string
// 	tags          []string
// }

func cloneRepo(cfg config.Git, url string) error {
	cmd := exec.Command("git", "clone", url)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(
		os.Environ(),
		useSSHKey(cfg.SSHKeyPath),
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	return nil
}

func useSSHKey(path string) string {
	return fmt.Sprintf("GIT_SSH_COMMAND=ssh -i %s -o IdentitiesOnly=yes", path)
}
