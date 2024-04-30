package commands

import "google.golang.org/api/youtube/v3"

type CommandConfig struct {
	YtClient  *youtube.Service
	VideoId   string
	HighestBy string
}
