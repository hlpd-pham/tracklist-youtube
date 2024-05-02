package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hlpd-pham/tracklist-youtube/commands"
	"github.com/hlpd-pham/tracklist-youtube/spotify_wrapper"
	"github.com/hlpd-pham/tracklist-youtube/yt"
)

func parseInput(text string) []string {
	lowercase := strings.ToLower(text)
	trimmed := strings.Trim(lowercase, " ")
	tokens := strings.Split(trimmed, " ")
	return tokens
}

func main() {
	logger := log.Default()
	tracklistCmd := flag.NewFlagSet("tracklist", flag.ExitOnError)
	videoIdFlag := tracklistCmd.String("videoId", "", "videoId")
	highestByFlag := tracklistCmd.String("highestBy", "", "highestBy")

	ytService, err := yt.NewYoutubeService()
	if err != nil {
		fmt.Println(fmt.Errorf("encounter error while creating youtube client: %v", err))
		os.Exit(1)
	}

	ytClient, err := yt.NewYtWrapperClient(logger, ytService)
	if err != nil {
		fmt.Println(fmt.Errorf("encounter error while creating youtube client: %v", err))
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		tracklistCmd.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "tracklist":
		tracklistCmd.Parse(os.Args[2:])
		logger.Println(*videoIdFlag, *highestByFlag)
		cmdCfg := commands.CommandConfig{
			YtClient:      ytClient,
			SpotifyClient: spotify_wrapper.GetWrapperClient(logger),
			VideoId:       *videoIdFlag,
			HighestBy:     *highestByFlag,
		}
		err = commands.CommandTracklist(&cmdCfg)
		if err != nil {
			tracklistCmd.Usage()
			os.Exit(1)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}
