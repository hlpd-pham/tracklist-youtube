package spotify_wrapper

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type WrapperClient struct {
	client *spotify.Client
	logger *log.Logger
}

func GetWrapperClient(logger *log.Logger) *WrapperClient {
	godotenv.Load()
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
		os.Exit(1)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	return &WrapperClient{client: spotify.New(httpClient), logger: logger}
}

func (c *WrapperClient) GetSongsFromLines(lines []string) {
	for _, line := range lines {
		// tokens := strings.Split(line, "-")
		bestEffortToken := strings.TrimSpace(line)
		if bestEffortToken == "ID" {
			continue
		}
		// fmt.Printf("bestEffortToken: '%s' - line: '%s'\n", bestEffortToken, line)

		ctx := context.Background()
		result, err := c.client.Search(ctx, bestEffortToken, spotify.SearchTypeTrack)
		if err != nil {
			fmt.Println(fmt.Errorf("found error while searching for song: %s, err: %s", bestEffortToken, err.Error()))
			continue
		}
		if len(result.Tracks.Tracks) == 0 {
			fmt.Printf("Did not find results for song: %s", bestEffortToken)
			continue
		}

		fmt.Printf("top result for song: %s\n", bestEffortToken)
		fmt.Println(result.Tracks.Tracks[0].ExternalURLs["spotify"])
		fmt.Println("--------------------------------")
	}
}
