package commands

import (
	"fmt"
	"log"
	"strings"
)

func CommandTracklist(config *CommandConfig) error {
	log.Printf("Getting tracklist for videoId '%s'\n", config.VideoId)
	config.YtClient.GetVideoInfo([]string{"snippet", "id"}, config.VideoId)

	tracklistComment, err := config.YtClient.GetTracklistComment(
		[]string{"snippet", "id", "replies"},
		config.VideoId, config.HighestBy)
	if err != nil {
		log.Printf("Error getting tracklist comment: %s\n", err)
		return err
	}
	fmt.Println(tracklistComment)
	fmt.Printf("config.FetchSpotify: %v\n", config.FetchSpotify)

	if config.FetchSpotify {
		log.Println("Fetching Spotify songs")
		lines := strings.Split(tracklistComment, "\n")
		trackResults := config.SpotifyClient.GetSongsFromLines(lines)
		for _, track := range *trackResults {
			fmt.Printf("Song: %s, Artist: %s\n", track.Name, track.Artists[0].Name)
		}
	}

	return nil
}
