package commands

import (
	"fmt"

	"github.com/hlpd-pham/tracklist-youtube/yt"
)

func CommandTracklist(config *CommandConfig) error {
	videoId := config.CommandArgs[0]
	fmt.Printf("Getting tracklist for videoId '%s'\n", videoId)
	yt.GetVideoInfo(config.YtClient, []string{"snippet", "id"}, videoId)
	yt.GetTracklistCommentByLength(config.YtClient, []string{"snippet", "id", "replies"}, videoId)

	return nil
}
