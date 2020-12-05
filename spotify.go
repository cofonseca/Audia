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
	BPM        int
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

func connectToSpotify(clientID string, clientSec string) (spotify.Client, string) {
	// Configure Client & Connect
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSec,
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("Failed to get auth token: %v", err)
	}
	client := spotify.Authenticator{}.NewClient(token)
	return client, token.AccessToken
}

func convertPlaylistURLtoID(playlistURL string) string {
	playlistID := strings.Split(playlistURL, "/playlist/")[1]
	playlistID = strings.Split(playlistID, "?")[0]
	return playlistID
}

func getTrackAttributes(token string, playlistID string) trackAttributes {
	URL := "https://api.spotify.com/v1/audio-features/" + playlistID
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

// TODO: PAGINATION!!! This only returns 100 items,
// even if there are more than 100 items in the playlist.
func getPlaylistContents(client spotify.Client, token string, playlistID string) playlist {
	playlistTracks, err := client.GetPlaylistTracks(spotify.ID(playlistID))
	if err != nil {
		fmt.Println("Failed to retreive playlist information:", err)
	}
	fmt.Println("Number of tracks in playlist:", playlistTracks.Total)
	fmt.Println("")

	var p playlist

	for _, t := range playlistTracks.Tracks {
		var track track
		track.Artist = t.Track.Artists[0].Name
		track.ID = t.Track.ID.String()
		track.Title = t.Track.Name
		// TODO: getTrackAttributes can be done in a single call for all tracks.
		// See Spotify API docs for a way to implement this.
		track.Attributes = getTrackAttributes(token, track.ID)
		track.BPM = int(math.RoundToEven(track.Attributes.Tempo))
		p.Tracks = append(p.Tracks, track)
	}

	return p
}
