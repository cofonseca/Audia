package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

type youtubeAPICredentials struct {
	APIKeys []string `json:"youtubeApiKeys"`
}

type youtubeVideo struct {
	Title string `json:"title"`
	ID    string `json:"encrypted_id"`
	URL   string
}

type youtubeSearchResults struct {
	Videos []youtubeVideo `json:"video"`
}

func connectToYoutubeByAPI(APIKeys []string) *youtube.Service {
	// Configure Client & Connect
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(APIKeys[1]))
	if err != nil {
		fmt.Println("Failed to create YouTube service:", err)
	}

	return service
}

func searchYoutubeForTrackByAJAX(conf config, track track) youtubeVideo {
	var video youtubeVideo
	var ytResults youtubeSearchResults
	youtubeSearchURL := "https://www.youtube.com/search_ajax?style=json&search_query="
	artist := strings.ReplaceAll(track.Artist, " ", "+")
	title := strings.ReplaceAll(track.Title, " ", "+")
	queryURL := youtubeSearchURL + fmt.Sprintf("%s+-+%s&page=0&hl=en", artist, title)

	client := &http.Client{}
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		fmt.Println("Error building YouTube AJAX request:", err)
	}
	req.Header.Add("x-youtube-client-name", conf.YoutubeClientName)
	req.Header.Add("x-youtube-client-version", conf.YoutubeClientVer)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to make AJAX request:", err)
	}
	defer resp.Body.Close()

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Couldn't unmarshal JSON response:", err)
	}
	json.Unmarshal(jsonBody, &ytResults)
	video = ytResults.Videos[0]
	youtubeBaseURL := "https://www.youtube.com/watch?v="
	video.URL = fmt.Sprintf("%s%s", youtubeBaseURL, video.ID)
	return video
}

func searchYoutubeForTrackByAPI(service *youtube.Service, track track) youtubeVideo {
	// Perform Search
	list := []string{"snippet", "id"}
	//list = append(list, "snippet", "id")
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
		fmt.Println("Video not found for track", track.Title)
	}
	return video
}

func ytAPIWorker(id int, jobs <-chan int, results chan<- int, service *youtube.Service, playlist playlist) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		result := searchYoutubeForTrackByAPI(service, playlist.Tracks[j-1])
		getAudioFromVideo(result, playlist.Tracks[j-1], input.Destination)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func ytAJAXWorker(id int, jobs <-chan int, results chan<- int, conf config, playlist playlist) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		result := searchYoutubeForTrackByAJAX(conf, playlist.Tracks[j-1])
		getAudioFromVideo(result, playlist.Tracks[j-1], input.Destination)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}
