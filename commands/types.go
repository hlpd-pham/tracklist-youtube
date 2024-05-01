package commands

import (
	"github.com/hlpd-pham/tracklist-youtube/spotify_wrapper"
	"google.golang.org/api/youtube/v3"
)

type CommandConfig struct {
	YtClient      *youtube.Service
	SpotifyClient *spotify_wrapper.WrapperClient
	VideoId       string
	HighestBy     string
}
