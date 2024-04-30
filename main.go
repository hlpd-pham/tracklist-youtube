package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hlpd-pham/tracklist-youtube/commands"
	"github.com/hlpd-pham/tracklist-youtube/yt"
)

func parseInput(text string) []string {
	lowercase := strings.ToLower(text)
	trimmed := strings.Trim(lowercase, " ")
	tokens := strings.Split(trimmed, " ")
	return tokens
}

func main() {
	tracklistCmd := flag.NewFlagSet("tracklist", flag.ExitOnError)
	videoIdFlag := tracklistCmd.String("videoId", "", "videoId")
	highestByFlag := tracklistCmd.String("highestBy", "like", "highestBy")

	service, err := yt.CreateYoutubeClient()
	if err != nil {
		fmt.Println(fmt.Errorf("encounter error while creating youtube client: %v", err))
		return
	}

	if len(os.Args) < 2 {
		tracklistCmd.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "tracklist":
		tracklistCmd.Parse(os.Args[2:])
		cmdCfg := commands.CommandConfig{
			YtClient:  service,
			VideoId:   *videoIdFlag,
			HighestBy: *highestByFlag,
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
