package main

import (
	"backend/structs"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func getNowPlaying(w http.ResponseWriter, nowPlaying *structs.NowPlaying) {
	url1 := fmt.Sprintf("https://api.listenbrainz.org/1/user/%s/playing-now", lbUser)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		log.Println("Error creating HTTP request: ", err)
		http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", lbToken))
	res, err := client.Do(req) // send the request
	if err != nil {
		log.Println("Error sending HTTP request: ", err)
		http.Error(w, "Error sending `Now Playing` request", http.StatusInternalServerError)
		return
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Println("Error reading HTTP response body: ", readErr)
		http.Error(w, "Error reading `Now Playing` response", http.StatusInternalServerError)
		return
	}

	svc1 := structs.LBPlayingNow{}
	jsonErr := json.Unmarshal(body, &svc1)
	if jsonErr != nil {
		log.Println("Error unmarshaling JSON response: ", jsonErr)
		http.Error(w, "Error parsing `Now Playing` response", http.StatusInternalServerError)
		return
	}

	if len(svc1.Payload.Listens) > 0 {
		nowPlaying.Artist = svc1.Payload.Listens[0].TrackMetadata.ArtistName
		nowPlaying.Album = svc1.Payload.Listens[0].TrackMetadata.ReleaseName
		nowPlaying.Track = svc1.Payload.Listens[0].TrackMetadata.TrackName
	} else {
		// no song currently playing
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func getAlbumMBID(w http.ResponseWriter, nowPlaying *structs.NowPlaying) (string, error) {
	fmt.Println("GET Album MBID...")
	client := &http.Client{}
	// TODO: this looks bad
	apiURL := "http://musicbrainz.org/ws/2/release/?query=artist:\"" + url.QueryEscape(nowPlaying.Artist) + "\"%20AND%20release:\"" + url.QueryEscape(nowPlaying.Album) + "\"&fmt=json&limit=5"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("Error creating HTTP request: ", err)
		http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
		return "", err
	}

	req.Header.Set("User-Agent", "AlbumartCollage/0.1 ( josh.l8n@gmail.com )")
	res, err := client.Do(req) // send the request
	if err != nil {
		log.Println("Error sending HTTP request: ", err)
		http.Error(w, "Error sending `Now Playing` request", http.StatusInternalServerError)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading HTTP response body: ", err)
		http.Error(w, "Error reading `Now Playing` response", http.StatusInternalServerError)
		return "", err
	}

	mbResponse := structs.MBAlbumInfo{}
	err = json.Unmarshal(body, &mbResponse)
	if err != nil {
		log.Println("Error unmarshaling JSON response: ", err, apiURL)
		http.Error(w, "Error parsing `Now Playing 2` response", http.StatusInternalServerError)
		return "", err
	}

	if len(mbResponse.Releases) > 0 {
		for _, release := range mbResponse.Releases {
			if strings.ToLower(release.Title) == strings.ToLower(nowPlaying.Album) &&
				strings.ToLower(release.ArtistCredit[0].Name) == strings.ToLower(nowPlaying.Artist) &&
				strings.ToLower(release.ReleaseGroup.PrimaryType) == "album" &&
				strings.ToLower(release.Disambiguation) == "" {

				return release.ID, nil
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return "", errors.New("Album not found while searching MusicBrainz")
}

func getCoverArtURL(w http.ResponseWriter, nowPlaying *structs.NowPlaying) (string, error) {
	fmt.Println("GET Coverart...", nowPlaying.AlbumMBID)
	client := &http.Client{} // TODO: should I include this in every function or pass it?
	apiURL := fmt.Sprintf("https://coverartarchive.org/release/%s", nowPlaying.AlbumMBID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("Error creating HTTP request: ", err)
		http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
		return "", err
	}

	req.Header.Set("User-Agent", "AlbumArtCollage/0.1 ( josh.l8n@gmail.com )")
	res, err := client.Do(req) // send the request
	if err != nil {
		log.Println("Error sending HTTP request: ", err)
		http.Error(w, "Error sending `Now Playing` request", http.StatusInternalServerError)
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("cover art not found: %d", res.StatusCode)
		fmt.Println(errMsg)
		return "", errors.New(errMsg)
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Println("Error reading HTTP response body: ", readErr)
		http.Error(w, "Error reading `Now Playing` response", http.StatusInternalServerError)
		return "", err
	}

	caResponse := structs.CoverartArchive{}
	jsonErr := json.Unmarshal(body, &caResponse)
	if jsonErr != nil {
		log.Println("Error unmarshaling JSON response: ", jsonErr, apiURL)
		http.Error(w, "Error parsing `Now Playing` response", http.StatusInternalServerError)
		return "", err
	}

	// TODO: fallback to smaller sizes if large not found
	return caResponse.Images[0].Image, nil
}
