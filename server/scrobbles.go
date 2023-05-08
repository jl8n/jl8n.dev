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

	sbResponse := structs.LBPlayingNow{}
	jsonErr := json.Unmarshal(body, &sbResponse)
	if jsonErr != nil {
		log.Println("Error unmarshaling JSON response: ", jsonErr)
		http.Error(w, "Error parsing `Now Playing` response", http.StatusInternalServerError)
		return
	}

	if len(sbResponse.Payload.Listens) > 0 {
		nowPlaying.Artist = sbResponse.Payload.Listens[0].TrackMetadata.ArtistName
		nowPlaying.Album = sbResponse.Payload.Listens[0].TrackMetadata.ReleaseName
		nowPlaying.Track = sbResponse.Payload.Listens[0].TrackMetadata.TrackName
	} else {
		// no song currently playing
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func getAlbumMBID(w http.ResponseWriter, nowPlaying *structs.NowPlaying) (string, error) {
	fmt.Println("GET Album MBID...")
	client := &http.Client{}

	// Construct URL
	query := fmt.Sprintf("artist:%s AND release:%s", nowPlaying.Artist, nowPlaying.Album)
	params := fmt.Sprintf("query=%s&fmt=json&limit=5", url.QueryEscape(query))
	apiURL := url.URL{
		Scheme:   "http",
		Host:     "musicbrainz.org",
		Path:     "ws/2/release/",
		RawQuery: params,
	}

	fmt.Println("mbURL:", apiURL.String())

	req, err := http.NewRequest("GET", apiURL.String(), nil)
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

	fmt.Println(4, apiURL)

	if len(mbResponse.Releases) > 0 {

		releasesCount := len(mbResponse.Releases)
		weights := make([]int, releasesCount)

		// create score for each
		for i, release := range mbResponse.Releases {
			// TODO: REGEX, better var name
			one, two := strings.TrimSpace(strings.ToLower(release.Title)), strings.TrimSpace(strings.ToLower(nowPlaying.Album))
			one, two = strings.ReplaceAll(one, ",", ""), strings.ReplaceAll(two, ",", "")
			one, two = strings.ReplaceAll(one, "'", ""), strings.ReplaceAll(two, "'", "")
			one, two = strings.ReplaceAll(one, "’", ""), strings.ReplaceAll(two, "’", "")

			isTitleSame := one == two
			isArtistSame := strings.ToLower(release.ArtistCredit[0].Name) == strings.ToLower(nowPlaying.Artist)
			isAlbum := strings.ToLower(release.ReleaseGroup.PrimaryType) == "album" ||
				strings.ToLower(release.ReleaseGroup.PrimaryType) == "ep" ||
				strings.ToLower(release.ReleaseGroup.PrimaryType) == "single"
			isNotDisambig := strings.ToLower(release.Disambiguation) == ""
			isOfficial := release.Status == "Official"
			//isDigital := release.Media[0].Format == "Digital Media"
			isPhysical := release.Media[0].Format == "Album"

			fmt.Println(isTitleSame, isArtistSame, isAlbum, isNotDisambig, isOfficial, isAlbum)

			for _, val := range []bool{isTitleSame, isArtistSame, isAlbum, isNotDisambig, isOfficial, isPhysical} {
				if val {
					weights[i]++
				}
			}

			// if isTitleSame && isArtistSame && isAlbum && isNotDisambig && isOfficial && isMedia {
			// 	fmt.Println("Best match:", release.ID)

			// 	// TODO: check if album-art can be fetched here
			// 	// if not, then move to the next album

			// }
		}

		winnerIndex := 0
		largest := 0

		for i, candidate := range weights {
			if candidate > largest {
				largest = candidate
				winnerIndex = i
			}
		}

		fmt.Println("winner", winnerIndex, mbResponse.Releases[winnerIndex].ID)

		return mbResponse.Releases[winnerIndex].ID, nil
		//return release.ID, nil

	}

	return "", errors.New("Suitable match not found while searching MusicBrainz")
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
