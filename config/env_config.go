package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	// EnvTeamsWebhook env var name for Teams webhook url
	EnvTeamsWebhook = "EXTRA_TEAMS_WEBHOOK"

	// EnvGithubToken env var name for Github token
	EnvGithubToken = "EXTRA_GITHUB_TOKEN"

	// EnvUpdateInterval env var name update interval in hours
	EnvUpdateInterval = "EXTRA_APP_INTERVAL"
)

func setEnvVarsInConfig(config *Config) {
	teamsWebHook := os.Getenv(EnvTeamsWebhook)
	if teamsWebHook != "" {
		config.Teams.WebhookURL = teamsWebHook
	}

	githubToken := os.Getenv(EnvGithubToken)
	if githubToken != "" {
		config.Github.AuthToken = githubToken
	}

	updateInterval := os.Getenv(EnvUpdateInterval)
	if updateInterval != "" {
		ui, err := strconv.Atoi(updateInterval)
		if err != nil {
			fmt.Printf("The envVar %s must be a number\n", EnvUpdateInterval)
		} else {
			config.App.CheckIntervalHours = ui
		}
	}

}
