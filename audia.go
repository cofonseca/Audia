package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	youtube "google.golang.org/api/youtube/v3"
)

type userInput struct {
	URL         string
	Destination string
	Workers     uint
}

var (
	input userInput
)

func validateDestinationPath(path string) {
	newPath := strings.ReplaceAll(path, "\\", "/")
	if string(newPath[len(newPath)-1:]) != "/" {
		newPath = newPath + "/"
	}
	input.Destination = newPath
}

func parseFlags() error {
	flag.StringVar(&input.URL, "url", "", "The full URL of the song or playlist.")
	flag.StringVar(&input.Destination, "destination", "", "The absolute path to a folder where the MP3 files should be saved.")
	flag.UintVar(&input.Workers, "workers", 1, "The number of songs to download at a time. A higher number will download more songs concurrently, meaning that a large playlist may download faster. A higher number will require more network bandwith and a more powerful CPU. General recommendation is 1 per logical processor.")
	flag.Parse()

	if input.URL == "" {
		return errors.New("Required parameter 'url' is missing. Please provide a URL by using -url <url>")
	}

	if input.Destination == "" {
		return errors.New("Required parameter 'destination' is missing. Please provide the destination path by using -destination")
	}

	validateDestinationPath(input.Destination)

	if input.Workers <= 0 || input.Workers >= 255 {
		return errors.New("Required parameter 'workers' is invalid. Please provide a number of workers by using -workers <size>. Workers must be a positive integer between 1 and 254")
	}

	return nil
}

func main() {
	err := parseFlags()
	if err != nil {
		fmt.Println(err)
		return
	}

	conf := parseConfig("./config.json")

	// Get Track List from Spotify
	spotifyClient, token := connectToSpotify(conf.SpotifyClientID, conf.SpotifyClientSec)
	playlistID := convertPlaylistURLtoID(input.URL)
	playlist := getPlaylistContents(spotifyClient, token, playlistID)

	numJobs := len(playlist.Tracks)
	fmt.Println(numJobs)
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Search YouTube for Track & Download
	var yt *youtube.Service
	if conf.UseYoutubeAPI == true {
		yt = connectToYoutubeByAPI(conf.YoutubeAPIKeys)

		for w := 1; w <= int(input.Workers); w++ {
			go ytAPIWorker(w, jobs, results, yt, playlist)
		}
	} else {
		for w := 1; w <= int(input.Workers); w++ {
			go ytAJAXWorker(w, jobs, results, conf, playlist)
		}
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		<-results
	}

}
