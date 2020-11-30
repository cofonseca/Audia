package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type spotifyAPICredentials struct {
	ClientID  string `json:"clientID"`
	ClientSec string `json:"clientSec"`
}

type track struct {
	Artist     string
	Title      string
	ID         string
	BPM        float64
	Attributes trackAttributes
}

type trackAttributes struct {
	Danceability     float32 `json:"danceability"`
	Energy           float32 `json:"energy"`
	Key              int     `json:"key"`
	Loudness         float32 `json:"loudness"`
	Mode             int     `json:"mode"`
	Speechinness     float32 `json:"speechiness"`
	Acousticness     float32 `json:"acousticness"`
	Instrumentalness float32 `json:"instrumentalness"`
	Liveness         float32 `json:"liveness"`
	Valence          float32 `json:"valence"`
	Tempo            float64 `json:"tempo"`
	DurationMS       int32   `json:"duration_ms"`
	TimeSignature    int     `json:"time_signature"`
}

type playlist struct {
	Tracks []track
}

func connectToSpotify() (spotify.Client, string) {
	// Get API Creds from Config & Create Config Object
	var creds spotifyAPICredentials
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("Error reading from config.json:", err)
	}
	err = json.Unmarshal(configFile, &creds)
	if err != nil {
		fmt.Println("Error unmarshaling config JSON:", err)
	}

	config := &clientcredentials.Config{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSec,
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("Failed to get auth token: %v", err)
	}
	fmt.Println("Token:", token.AccessToken)
	client := spotify.Authenticator{}.NewClient(token)
	return client, token.AccessToken
}

func convertPlaylistURLtoID(playlistURL string) string {
	fmt.Println("Original URL:", playlistURL)
	spotifyID := strings.Split(playlistURL, "/playlist/")[1]
	spotifyID = strings.Split(spotifyID, "?")[0]
	fmt.Println("Stripped ID:", spotifyID)
	return spotifyID
}

func getTrackAttributes(token string, spotifyID string) trackAttributes {
	URL := "https://api.spotify.com/v1/audio-features/" + spotifyID
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Println("Error building getTrackBPM request:", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to make request:", err)
	}
	defer resp.Body.Close()

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Couldn't unmarshal JSON response:", err)
	}
	var attributes trackAttributes
	json.Unmarshal(jsonBody, &attributes)

	return attributes

}

func getPlaylistContents(client spotify.Client, token string, spotifyID string) {
	playlistTracks, err := client.GetPlaylistTracks(spotify.ID(spotifyID))
	if err != nil {
		fmt.Println("Failed to retreive playlist information:", err)
	}
	fmt.Println("Number of tracks:", playlistTracks.Total)
	fmt.Println("")

	for _, t := range playlistTracks.Tracks {
		var track track
		track.Artist = t.Track.Artists[0].Name
		track.ID = t.Track.ID.String()
		track.Title = t.Track.Name
		track.Attributes = getTrackAttributes(token, track.ID)
		track.BPM = math.RoundToEven(track.Attributes.Tempo)
		fmt.Println("")
		fmt.Println(track)
	}
	// Get each track and create a track struct
	// For each struct, get the BPM
	// Store all tracks in a playlist struct
}
