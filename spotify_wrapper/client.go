package spotify_wrapper

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/hlpd-pham/tracklist-youtube/util"
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

func (c *WrapperClient) GetSongsFromLines(lines []string) *[]spotify.FullTrack {
	trackResults := make([]spotify.FullTrack, 0)
	for _, line := range lines {
		bestEffortToken := strings.TrimSpace(line)
		if strings.Contains(strings.ToLower(bestEffortToken), "id") ||
			strings.Contains(strings.ToLower(bestEffortToken), "tracklist") {
			continue
		}

		if strings.Contains(bestEffortToken, "-") {
			allTokens := strings.Split(bestEffortToken, "-")
			bestEffortToken = strings.ToLower(strings.TrimSpace(allTokens[len(allTokens)-1]))
		}

		ctx := context.Background()
		result, err := c.client.Search(ctx, bestEffortToken, spotify.SearchTypeTrack)
		if err != nil {
			continue
		}
		if len(result.Tracks.Tracks) == 0 {
			c.logger.Printf("Did not find results for song: %s", bestEffortToken)
			continue
		}

		bestResult := c.findBestResult(bestEffortToken, result.Tracks.Tracks)
		artistNames := []string{}
		for _, aName := range bestResult.Artists {
			artistNames = append(artistNames, aName.Name)
		}

		trackResults = append(trackResults, *bestResult)
		c.logger.Printf("top result for song '%s' is %s by %v\n", bestEffortToken, bestResult.Name, artistNames)
	}
	return &trackResults
}

func (c *WrapperClient) findBestResult(query string, searchResult []spotify.FullTrack) *spotify.FullTrack {
	var bestTrack *spotify.FullTrack
	bestScore := 0

	for _, track := range searchResult {
		trackScore, isFullyMatched := c.calculateSearchScore(query, track)
		if isFullyMatched {
			return &track
		}
		if trackScore > bestScore {
			bestTrack = &track
		}
	}
	if bestTrack == nil {
		return &searchResult[0]
	}
	return bestTrack
}

// return score of a track and if the track is fully matched with query
func (c *WrapperClient) calculateSearchScore(query string, track spotify.FullTrack) (int, bool) {
	queryTokens := strings.Split(query, " ")
	songNameTokens := strings.Split(track.Name, " ")
	for index, sToken := range songNameTokens {
		songNameTokens[index] = strings.ToLower(sToken)
	}

	artistTokens := []string{}
	for _, artist := range track.Artists {
		curArtistTokens := strings.Split(artist.Name, " ")
		for index, cToken := range curArtistTokens {
			curArtistTokens[index] = strings.ToLower(cToken)
		}
		artistTokens = append(artistTokens, curArtistTokens...)
	}

	artistScore := 0
	matchIndexes := []int{}
	for index, qToken := range queryTokens {
		if util.Contains(artistTokens, qToken) {
			artistScore += 2
			matchIndexes = append(matchIndexes, index)
		}
	}
	newQueryTokens := []string{}
	for index, qToken := range queryTokens {
		if !util.Contains(matchIndexes, index) {
			newQueryTokens = append(newQueryTokens, qToken)
		}
	}
	queryTokens = newQueryTokens

	songScore := 0
	for _, qToken := range queryTokens {
		if util.Contains(songNameTokens, qToken) {
			songScore += 1
		}
	}

	isFullyMatched := artistScore == len(artistTokens) && songScore == len(songNameTokens)
	// c.logger.Printf("track: %v, artists: %v, score: %d, isFullyMatch: %v\n",
	// 	songNameTokens, artistTokens, songScore+artistScore, isFullyMatched)

	if isFullyMatched {
		return artistScore + songScore, true
	}

	return artistScore + songScore, false
}
