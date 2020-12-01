package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func getAudioFromVideo(video youtubeVideo, track track, destDir string) {
	destination := fmt.Sprintf("%s%d %s - %s.%s", destDir, track.BPM, track.Artist, track.Title, "%(ext)s")
	cmd := exec.Command("youtube-dl", "-x", "--audio-format", "mp3", "-o", destination, video.URL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	outLines := strings.Split(outStr, "\n")
	fmt.Println(outLines[len(outLines)-3])
	if errStr != "" {
		fmt.Println(errStr)
	}
}
