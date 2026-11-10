package main

import (
	"argus-agent/cmd"
	"argus-agent/internal/config"
	"argus-agent/internal/git"
	"argus-agent/internal/llm"
	"argus-agent/internal/llm/openai"
	"argus-agent/internal/report"
	"fmt"
	"os"
)

var (
	appCfg      *config.Application
	llmProvider llm.Provider
	defaultDir  string
)

func init() {
	var err error
	appCfg, err = config.Load()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	llmProvider = openai.NewProvider(openai.NewClient(appCfg.LLM))

	defaultDir, err = report.DefaultDir()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	reportService, err := report.NewService(defaultDir)
	if err != nil {
		fmt.Println("failed to create report service: ", err)
		os.Exit(1)
	}

	gitSvc, err := git.NewService(appCfg.Git, llmProvider, reportService)
	if err != nil {
		fmt.Println("failed to create git service: ", err)
		os.Exit(1)
	}

	cmd.InitCommands(gitSvc, config.EditorService)
	cmd.Execute()
}
