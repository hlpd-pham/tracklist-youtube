package yt

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

func CreateYoutubeClient() (*youtube.Service, error) {
	ctx := context.Background()

	b, err := os.ReadFile("youtube_client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/tracklist-youtube.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeForceSslScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {

		HandleError(err, "Error creating YouTube client")
		return nil, err
	}
	return service, nil
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = GetTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// GetTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0o700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("tracklist-youtube.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func HandleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func replaceHtmlEntities(input string) string {
	// Define a regular expression pattern to match HTML entities
	pattern := `(&amp;)|(&#39;)`

	// Compile the regular expression pattern
	re := regexp.MustCompile(pattern)

	// Replace HTML entities with their decoded counterparts
	return re.ReplaceAllStringFunc(input, func(match string) string {
		return html.UnescapeString(match)
	})
}

func GetTracklistComment(service *youtube.Service,
	parts []string,
	videoId string,
	highestBy string,
) (string, error) {
	call := service.CommentThreads.List(parts).VideoId(videoId)
	response, err := call.Do()
	HandleError(err, "")

	trackListComments := []*youtube.Comment{}

	for response.NextPageToken != "" {
		for _, commentThread := range response.Items {
			topComment := commentThread.
				Snippet.
				TopLevelComment
			if strings.Contains(topComment.Snippet.TextDisplay, "Tracklist") ||
				strings.Contains(topComment.Snippet.TextDisplay, "tracklist") {
				trackListComments = append(trackListComments, topComment)
			}
		}
		call = service.CommentThreads.
			List(parts).
			VideoId(videoId).
			PageToken(response.NextPageToken)
		response, err = call.Do()
		HandleError(err, "")
	}

	bestComment := youtube.Comment{Snippet: &youtube.CommentSnippet{TextDisplay: "", LikeCount: 0}}
	if len(trackListComments) == 0 {
		return "", errors.New("no tracklist comment found")
	}

	fmt.Printf("Found %d comments\n", len(trackListComments))
	for _, comment := range trackListComments {
		if highestBy == "length" {
			if len(comment.Snippet.TextDisplay) > len(bestComment.Snippet.TextDisplay) {
				bestComment = *comment
			}
		} else {
			if comment.Snippet.LikeCount > bestComment.Snippet.LikeCount {
				bestComment = *comment
			}
		}
	}
	fmt.Printf("Best comment has %d likes and length %d\n",
		bestComment.Snippet.LikeCount,
		len(bestComment.Snippet.TextDisplay))
	pattern := `<a.*>.*<\/a>\s`
	re := regexp.MustCompile(pattern)
	removeBreakTags := strings.ReplaceAll(
		bestComment.Snippet.TextDisplay, "<br>", "\n")
	result := re.ReplaceAllString(removeBreakTags, "")

	return replaceHtmlEntities(result), nil
}

func GetVideoInfo(service *youtube.Service, parts []string, videoId string) error {
	call := service.Videos.List(parts).Id(videoId)
	response, err := call.Do()
	HandleError(err, "")
	if len(response.Items) == 0 {
		return fmt.Errorf("could not find video for id: %s", videoId)
	}
	fmt.Printf("Video Id: %s, title: %s\n", videoId, response.Items[0].Snippet.Title)
	return nil
}

func ChannelsListByHandle(service *youtube.Service, part string, forUsername string) {
	channels := []string{part}
	call := service.Channels.List(channels)
	call = call.ForHandle(forUsername)
	response, err := call.Do()
	HandleError(err, "")
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
}
