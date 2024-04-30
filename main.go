package main

import (
	"fmt"

	"github.com/hlpd-pham/tracklist-youtube/yt"
)

func main() {
	service, err := yt.CreateYoutubeClient()
	if err != nil {
		fmt.Println(fmt.Errorf("encounter error while creating youtube client: %v", err))
		return
	}

	videoId := "6V2kMynnQ7M"
	yt.GetVideoInfo(service, []string{"snippet", "id"}, videoId)
	yt.GetTracklistComment(service, []string{"snippet", "id", "replies"}, videoId)
}
