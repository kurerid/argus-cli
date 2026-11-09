// Package config
package config

type Application struct {
	Git
}

type Git struct {
	UserEmail  string
	UserName   string
	SSHKeyPath string
}

func Load() (*Application, error) {
	app := Application{}
	gitCfg, err := loadGitConfig()
	if err != nil {
		return nil, err
	}
	app.Git = gitCfg
	return &app, nil
}

func loadGitConfig() (Git, error) {
	return Git{
		UserEmail:  "kurerid@yandex.ru",
		UserName:   "vasnin10a22syu",
		SSHKeyPath: "~/.ssh/id_ed25519.pub",
	}, nil
}
