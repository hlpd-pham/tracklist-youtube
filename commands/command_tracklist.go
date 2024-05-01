package commands

import (
	"fmt"
	"strings"
)

func CommandTracklist(config *CommandConfig) error {
	fmt.Printf("Getting tracklist for videoId '%s'\n", config.VideoId)
	config.YtClient.GetVideoInfo([]string{"snippet", "id"}, config.VideoId)

	tracklistComment, err := config.YtClient.GetTracklistComment(
		[]string{"snippet", "id", "replies"},
		config.VideoId, config.HighestBy)
	if err != nil {
		return err
	}
	lines := strings.Split(tracklistComment, "\n")
	config.SpotifyClient.GetSongsFromLines(lines)

	return nil
}
