package commands

import "google.golang.org/api/youtube/v3"

type CommandConfig struct {
	YtClient    *youtube.Service
	CommandArgs []string
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(cfg *CommandConfig) error
}
