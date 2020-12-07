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

type audioFeaturesList struct {
	AttributesList []trackAttributes `json:"audio_features"`
}

type trackAttributeMap struct {
	AttributeMap map[string]trackAttributes
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
	TrackID          string  `json:"id"`
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

func getTrackBPM(maps []trackAttributeMap, pl playlist) playlist {
	var p playlist
	for _, track := range pl.Tracks {
		for _, m := range maps {
			if _, exists := m.AttributeMap[track.ID]; exists {
				track.Attributes = m.AttributeMap[track.ID]
				track.BPM = int(math.RoundToEven(m.AttributeMap[track.ID].Tempo))
				p.Tracks = append(p.Tracks, track)
			}

		}
	}

	return p
}

func getAttributesOfTracks(token string, tracks []track) trackAttributeMap {
	URL := "https://api.spotify.com/v1/audio-features/?ids="
	client := &http.Client{}

	var IDList string
	for _, track := range tracks {
		IDList = fmt.Sprintf("%s%s,", IDList, track.ID)
	}
	IDList = strings.TrimSuffix(IDList, ",")

	queryURL := fmt.Sprintf("%s%s", URL, IDList)

	req, err := http.NewRequest("GET", queryURL, nil)
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
	var attributesList audioFeaturesList
	json.Unmarshal(jsonBody, &attributesList)

	var attributeMap trackAttributeMap
	attributeMap.AttributeMap = make(map[string]trackAttributes)
	for _, a := range attributesList.AttributesList {
		attributeMap.AttributeMap[a.TrackID] = a
	}

	return attributeMap

}

func getPlaylistContents(client spotify.Client, token string, playlistID string) playlist {
	playlistTracks, err := client.GetPlaylistTracks(spotify.ID(playlistID))
	if err != nil {
		fmt.Println("Failed to retreive playlist information:", err)
	}

	var p playlist
	var attributes []trackAttributeMap

	for page := 1; ; page++ {
		var getAttributesList []track

		for _, t := range playlistTracks.Tracks {
			trk := track{
				Artist: t.Track.Artists[0].Name,
				ID:     t.Track.ID.String(),
				Title:  t.Track.Name,
			}
			p.Tracks = append(p.Tracks, trk)
			getAttributesList = append(getAttributesList, trk)
		}

		attributes = append(attributes, getAttributesOfTracks(token, getAttributesList))

		err = client.NextPage(playlistTracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	pl := getTrackBPM(attributes, p)
	return pl
}
