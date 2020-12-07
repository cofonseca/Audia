# Audia - Save Your Playlists!
Audia is a simple command-line application that enables you to download all songs in a given playlist.

For example, a DJ would need to have their .mp3 files stored locally in order to mix music, but organizing music is difficult, and downloading songs individually is an extremely slow and tedious process. Using Audia, you can simply provide the URL to a Spotify playlist, and Audia will automatically find all of the songs, download them to a given location on your computer, and rename them in the format of "BPM - Artist - Title".

## Prerequisites
Audia uses two popular tools under the hood: youtube-dl to download videos, and ffmpeg to extract audio. Both of these tools must be installed, and the path to each executable should be in your PATH.

Alternatively, if you don't want to install and set up dependencies, a Docker image is available and can easily be executed on any operating system.

## Usage
### Parameters
Audia requires three parameters:
- URL: The full Spotify URL of the playlist that you'd like to download. To find this, open Spotify, right-click on your playlist, go to Share, and click on Copy Playlist Link.
- Destination: The full path to a folder on your hard drive where you would like your music to be saved.
- Workers: The number of songs to download at a time, between 1 and 254. A higher number will download more songs at the same time, but will require a more powerful CPU and more network bandwidth. Recommendation is 1 per logical processor.

### Binary
```
audia.exe -url <URL> -destination <PATH> -workers <NUMBER>
```
Example:
```
audia.exe -url https://open.spotify.com/playlist/37i9dQZF1DX4dyzvuaRJ0n?si=-zG9PaXMReO2vVJ-YXvncA -destination C:\Users\jsmith\Music -workers 8
```
### Docker
```
docker run -it --rm --name audia -v <PATH>:/out -e WORKERS=<NUMBER> -e URL=<URL> gcr.io/rebred/audia:latest
```
Example:
```
docker run -it -rm --name audia -v C:/Users/jsmith/Music:/out -e WORKERS=8 -e URL=https://open.spotify.com/playlist/37i9dQZF1DX4dyzvuaRJ0n?si=-zG9PaXMReO2vVJ-YXvncA gcr.io/rebred-296012/audia:latest
```


## To Do:
- Add SoundCloud Support to match Spotify.
- Increase YouTube API key quota, add more keys, or get info from YouTube some other way not requiring an API key.
- Spotify results are paginated. A playlist might have >100 songs, but we only grab info for the first 100. This should be fixed.