package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type spotifyAPICredentials struct {
	ClientID  string `json:"clientID"`
	ClientSec string `json:"clientSec"`
}

type track struct {
	Artist string
	Title  string
	ID     string
	BPM    int
}

type playlist struct {
	Tracks []track
}

func connectToSpotify() spotify.Client {
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

	client := spotify.Authenticator{}.NewClient(token)
	return client
}

func convertSpotifyURLtoID(playlistURL string) string {
	fmt.Println("Original URL:", playlistURL)
	spotifyID := strings.Split(playlistURL, "/playlist/")[1]
	spotifyID = strings.Split(spotifyID, "?")[0]
	fmt.Println("Stripped ID:", spotifyID)
	return spotifyID
}

func getPlaylistContents(client spotify.Client, spotifyID string) {
	playlistTracks, err := client.GetPlaylistTracks(spotify.ID(spotifyID))
	if err != nil {
		fmt.Println("Failed to retreive playlist information:", err)
	}
	fmt.Println("Number of tracks:", playlistTracks.Total)
	fmt.Println("")
	fmt.Println(playlistTracks.Tracks[0].Track.Artists)
	// Get each track and create a track struct
	// For each struct, get the BPM
	// Store all tracks in a playlist struct
}
