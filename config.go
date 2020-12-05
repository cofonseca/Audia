package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type config struct {
	SpotifyClientID   string   `json:"clientID"`
	SpotifyClientSec  string   `json:"clientSec"`
	UseYoutubeAPIKeys bool     `json:"useYoutubeApiKeys"`
	YoutubeAPIKeys    []string `json:"youtubeApiKeys"`
}

func parseConfig(filepath string) config {
	var conf config
	configFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading from config.json:", err)
	}
	err = json.Unmarshal(configFile, &conf)
	if err != nil {
		fmt.Println("Error unmarshaling config JSON:", err)
	}

	return conf
}
