package main

import (
	"argus-agent/cmd"
	"argus-agent/internal/config"
	"argus-agent/internal/git"
	"argus-agent/internal/llm/openai"
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

	openaiProvider := openai.NewProvider(openai.NewClient(appCfg.LLM))

	reportService, err := report.NewService()
	if err != nil {
		fmt.Println("failed to create report service: ", err)
		os.Exit(1)
	}

	gitSvc, err := git.NewService(appCfg.Git, openaiProvider, reportService)
	if err != nil {
		fmt.Println("failed to create git service: ", err)
		os.Exit(1)
	}

	cmd.InitCommands(gitSvc)
	cmd.Execute()
}
