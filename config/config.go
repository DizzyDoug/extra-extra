package config

var configuration Config

// Config contains all available configurations
type Config struct {
	Github GithubConfig
	App    AppConfig
	Teams  TeamsConfig
}

// GithubConfig contains all configurations for Github sources
type GithubConfig struct {
	AuthToken string
}

// AppConfig contains all configurations for the app
type AppConfig struct {
	CheckIntervalHours int
}

// TeamsConfig contains all configurations for MS Teams
type TeamsConfig struct {
	WebhookURL string
}

func init() {
	c := Config{}
	// load config file

	// overwrite config with set env vars
	setEnvVarsInConfig(&c)

	configuration = c
}

// GetConfig returns loaded configuration
func GetConfig() Config {
	return configuration
}
