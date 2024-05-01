package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hlpd-pham/tracklist-youtube/yt"
)

func CommandTracklist(config *CommandConfig) error {
	if config.VideoId == "" {
		return errors.New("videoId must be provided")
	}
	fmt.Printf("Getting tracklist for videoId '%s'\n", config.VideoId)
	yt.GetVideoInfo(config.YtClient, []string{"snippet", "id"}, config.VideoId)

	switch config.HighestBy {
	case "like":
		tracklistComment, err := yt.GetTracklistComment(
			config.YtClient,
			[]string{"snippet", "id", "replies"},
			config.VideoId, "like")
		if err != nil {
			return err
		}
		lines := strings.Split(tracklistComment, "\n")
		config.SpotifyClient.GetSongsFromLines(lines)
	case "length":
		tracklistComment, err := yt.GetTracklistComment(
			config.YtClient,
			[]string{"snippet", "id", "replies"},
			config.VideoId, "length")
		if err != nil {
			return err
		}
		lines := strings.Split(tracklistComment, "\n")
		config.SpotifyClient.GetSongsFromLines(lines)
	default:
		return errors.New("must sort comment by 'like' or 'length'")
	}

	return nil
}
