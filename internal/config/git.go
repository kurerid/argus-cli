package config

import (
	"fmt"
	"os/exec"
	"strings"
)

type Git struct {
	UserEmail  string `json:"user_email"`
	UserName   string `json:"user_name"`
	SSHKeyPath string `json:"ssh_key_path"`
}

func loadSystemGitSettings() Git {
	email := gitGlobal("user.email")
	name := gitGlobal("user.name")
	return Git{
		UserEmail: email,
		UserName:  name,
	}
}

func gitGlobal(key string) string {
	out, err := exec.Command("git", "config", "--global", "--get", key).Output()
	if err != nil {
		fmt.Println("failed to get git config global ", key)
		return ""
	}
	return strings.TrimSpace(string(out))
}
