package main

import (
	"bufio"
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
	service, err := yt.CreateYoutubeClient()
	if err != nil {
		fmt.Println(fmt.Errorf("encounter error while creating youtube client: %v", err))
		return
	}

	cmdCfg := commands.CommandConfig{
		YtClient: service,
	}
	for {
		fmt.Print("tracklist-YT > ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		inputTokens := parseInput(scanner.Text())
		command, ok := commands.GetCommands()[inputTokens[0]]
		if !ok {
			fmt.Println("Unknown command")
			commands.CommandHelp(&cmdCfg)
		} else {
			cmdCfg.CommandArgs = inputTokens[1:]
			err := command.Callback(&cmdCfg)
			if err != nil {
				fmt.Printf("Found error while running command %s: %v\n", command.Name, err)
			}
		}

		fmt.Println()

	}
}
