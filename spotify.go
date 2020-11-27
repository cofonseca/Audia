package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type spotifyAPICredentials struct {
	ClientID  string `json:"clientID"`
	ClientSec string `json:"clientSec"`
}

type spotifyResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expiration  int    `json:"expires_in"`
}

func connectToSpotify() string {
	// Get Spotify API Credentials from Config
	reqURL := "https://accounts.spotify.com/api/token"
	var creds spotifyAPICredentials
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("Error reading from config.json:", err)
	}
	err = json.Unmarshal(configFile, &creds)
	if err != nil {
		fmt.Println("Error unmarshaling config JSON:", err)
	}

	// Construct the POST Request
	auth := base64.StdEncoding.EncodeToString([]byte(creds.ClientID + ":" + creds.ClientSec))
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Error building Spotify request:", err)
	}
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	// Send POST Request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer resp.Body.Close()

	// Unmarshal Response & Return Access Token
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
	var response spotifyResponse
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		fmt.Println("Error unmarshaling JSON response from Spotify:", err)
		return ""
	}

	return response.AccessToken
}

func getPlaylistContents() {

}
