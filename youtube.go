package main

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

type youtubeAPICredentials struct {
	APIKeys []string `json:"youtubeApiKeys"`
}

type youtubeVideo struct {
	Title string
	ID    string
	URL   string
}

type youtubeSearchResults struct {
	Videos []youtubeVideo
}

var ytResults youtubeSearchResults

func connectToYoutube(APIKeys []string) *youtube.Service {
	// Configure Client & Connect
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(APIKeys[1]))
	if err != nil {
		fmt.Println("Failed to create YouTube service:", err)
	}

	return service
}

func searchYoutubeForTrack(service *youtube.Service, track track) youtubeVideo {
	// Perform Search
	var list []string
	list = append(list, "snippet", "id")
	query := fmt.Sprintf("%s - %s", track.Artist, track.Title)
	search := service.Search.List(list).MaxResults(1).Q(query)
	result, err := search.Do()
	if err != nil {
		fmt.Println("Error performing YouTube search:", err)
	}

	// Build Return Object
	var video youtubeVideo
	youtubeBaseURL := "https://www.youtube.com/watch?v="
	if result.Items[0].Id.Kind == "youtube#video" {
		video.Title = result.Items[0].Snippet.Title
		video.ID = result.Items[0].Id.VideoId
		video.URL = youtubeBaseURL + result.Items[0].Id.VideoId
	} else {
		fmt.Println("Video not fount for track", track.Title)
	}
	return video
}
