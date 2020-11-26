package main

import (
	"errors"
	"flag"
	"fmt"
)

type userInput struct {
	URL         string
	Destination string
}

var input userInput

func parseFlags() error {
	flag.StringVar(&input.URL, "url", "", "The full URL of the song or playlist")
	flag.StringVar(&input.Destination, "destination", "", "The absolute path to a folder where the MP3 files should be saved")
	flag.Parse()

	if input.URL == "" {
		return errors.New("Required parameter 'url' is missing. Please provide a URL by using -url <url>")
	}

	if input.Destination == "" {
		return errors.New("Required parameter 'destination' is missing. Please provide the destination path by using -destination")
	}

	return nil
}

/*
audia download song --url https://open.spotify.com/track/jiofewjifoew --destination c:\users\cfonseca\music
audia download playlist --url https://open.spotify.com/playlist/jioewjfieow --destination c:\users\cfonseca\music
audia --help
*/

// Figure out how to take in required parameters
// Figure out how to display help
// Connect to Spotify API
// Get song info: Artist, Title, BPM
// Connect to YouTube API. We'll need more than 1 API key bc only 100 searches allowed per day per key.
// If search with API key results in an error, reconnect with the second API key
// Default to scraping the UI if we run out of keys
// Search YouTube for Artist - Title (exclude "music video") and grab the URL of the top result
// Figure out the best way to convert it into audio-only.
// Either use one of the online yt-to-mp3 converters or an ffmpeg wrapper or something?
// Convert the file, download it, and rename it to BPM Artist - Title.mp3, then put in --destination dir
// Could prob use goroutines for the searching, conversion, and downloading to speed things up

func main() {
	err := parseFlags()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(input)
}
