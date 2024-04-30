package commands

import (
	"errors"
	"fmt"

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
		yt.GetTracklistCommentByLike(config.YtClient, []string{"snippet", "id", "replies"}, config.VideoId)
	case "length":
		yt.GetTracklistCommentByLength(config.YtClient, []string{"snippet", "id", "replies"}, config.VideoId)
	default:
		return errors.New("must sort comment by 'like' or 'length'")
	}

	return nil
}
