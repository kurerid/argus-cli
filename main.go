package main

import (
	"argus-agent/cmd"
	"argus-agent/internal/config"
	"argus-agent/internal/git"
	"argus-agent/internal/report"
	"fmt"
	"os"
)

func main() {
	appCfg, err := config.Load()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	reportService, err := report.NewService()
	if err != nil {
		fmt.Println("failed to create report service: ", err)
		os.Exit(1)
	}
	gitSvc, err := git.NewService(
		appCfg.Git,
		"api key here",
		reportService,
	)
	if err != nil {
		fmt.Println("failed to create git service: ", err)
		os.Exit(1)
	}
	cmd.InitCommands(gitSvc)
	cmd.Execute()
}
