package main

import (
	"backend/structs"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func sendListenBrainzReq() (*http.Response, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.listenbrainz.org/1/user/%s/playing-now", lbUser)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating HTTP request: ", err)
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", lbToken))
	res, err := client.Do(req) // send the request
	if err != nil {
		log.Println("Error sending HTTP request: ", err)
		return nil, err
	}

	return res, nil
}

func parseListenBrainzRes(res *http.Response) (structs.NowPlaying, error) {
	body, readErr := io.ReadAll(res.Body)
	nowPlaying := structs.NowPlaying{}

	if readErr != nil {
		log.Println("Error reading HTTP response body: ", readErr)
		return nowPlaying, readErr
	}

	sbResponse := structs.LBPlayingNow{}
	jsonErr := json.Unmarshal(body, &sbResponse)
	if jsonErr != nil {
		log.Println("Error unmarshaling JSON response: ", jsonErr)
		return nowPlaying, jsonErr
	}

	if len(sbResponse.Payload.Listens) > 0 {
		// song is playing
		nowPlaying.Artist = sbResponse.Payload.Listens[0].TrackMetadata.ArtistName
		nowPlaying.Album = sbResponse.Payload.Listens[0].TrackMetadata.ReleaseName
		nowPlaying.Track = sbResponse.Payload.Listens[0].TrackMetadata.TrackName
		return nowPlaying, nil
	} else {
		// no song currently playing
		return nowPlaying, nil
	}
}

func getCurrentlyPlayingSong() (structs.NowPlaying, error) {
	var nowPlaying structs.NowPlaying

	res, err := sendListenBrainzReq()
	if err != nil {
		return nowPlaying, err
	}

	nowPlaying, err = parseListenBrainzRes(res)
	if err != nil {
		return nowPlaying, err
	}

	return nowPlaying, nil
}
