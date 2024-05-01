package commands

import (
	"github.com/hlpd-pham/tracklist-youtube/spotify_wrapper"
	"github.com/hlpd-pham/tracklist-youtube/yt"
)

type CommandConfig struct {
	YtClient      *yt.YtWrapperClient
	SpotifyClient *spotify_wrapper.WrapperClient
	VideoId       string
	HighestBy     string
}
