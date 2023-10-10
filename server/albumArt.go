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
	"os"
	"strings"
)

// TODO: check if art exists before downloading
func downloadAlbumArt(url *string, mbid *string) error {
	response, err := http.Get(*url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	file, err := os.Create(fmt.Sprintf("./album-art/%s.jpg", *mbid))
	if err != nil {
		//panic(err)
		return err

	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		//panic(err)

	}

	return nil
}

func getAlbumArtURL(mbid string) (string, error) {
	fmt.Println("fire download2")
	client := &http.Client{}
	apiURL := fmt.Sprintf("https://coverartarchive.org/release/%s", mbid)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("Error creating HTTP request: ", err)
		return "", err
	}

	req.Header.Set("User-Agent", "AlbumArtCollage/0.1 ( josh.l8n@gmail.com )")
	res, err := client.Do(req) // send the request
	if err != nil {
		log.Println("Error sending HTTP request: ", err)
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err := fmt.Sprintf("Cover art not found: %d", res.StatusCode)
		return "", errors.New(err)
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Println("Error reading HTTP response body: ", readErr)
		return "", err
	}

	caResponse := structs.CoverartArchive{}
	jsonErr := json.Unmarshal(body, &caResponse)
	if jsonErr != nil {
		log.Println("Error unmarshaling JSON response: ", jsonErr, apiURL)
		//http.Error(w, "Error parsing `Now Playing` response", http.StatusInternalServerError)
		return "", err
	}

	// TODO: fallback to smaller sizes if large not found
	return caResponse.Images[0].Image, nil
}

func getBestFromSearchResults(mbResponse *structs.MBAlbumInfo, nowPlaying *structs.NowPlaying) string {
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
		isArtistSame := strings.EqualFold(release.ArtistCredit[0].Name, nowPlaying.Artist)
		isAlbum := strings.ToLower(release.ReleaseGroup.PrimaryType) == "album" ||
			strings.ToLower(release.ReleaseGroup.PrimaryType) == "ep" ||
			strings.ToLower(release.ReleaseGroup.PrimaryType) == "single"
		isNotDisambig := strings.EqualFold(release.Disambiguation, "")
		isOfficial := release.Status == "Official"
		//isDigital := release.Media[0].Format == "Digital Media"
		isPhysical := strings.EqualFold(release.Media[0].Format, "Album")

		//fmt.Println(isTitleSame, isArtistSame, isAlbum, isNotDisambig, isOfficial, isAlbum)

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

	return mbResponse.Releases[winnerIndex].ID

}

// Use a song+artist+album combo to find a matching release ID on MusicBrainz
// Returns a MusicBrainz MBID with the matching album
func findAlbumMatch(nowPlaying *structs.NowPlaying) (string, error) {
	client := &http.Client{}
	mbResponse := structs.MBAlbumInfo{} // MusicBrainz API response

	// Construct URL
	query := fmt.Sprintf("artist:%s AND release:%s", nowPlaying.Artist, nowPlaying.Album)
	params := fmt.Sprintf("query=%s&fmt=json&limit=5", url.QueryEscape(query))
	apiURL := url.URL{
		Scheme:   "http",
		Host:     "musicbrainz.org",
		Path:     "ws/2/release/",
		RawQuery: params,
	}

	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		log.Println("Error creating HTTP request: ", err)
		return "", err
	}

	req.Header.Set("User-Agent", "AlbumartCollage/0.1 ( josh.l8n@gmail.com )")
	res, err := client.Do(req) // send the request
	if err != nil {
		log.Println("Error sending HTTP request: ", err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading HTTP response body: ", err)
		return "", err
	}

	if err := json.Unmarshal(body, &mbResponse); err != nil {
		log.Printf("Error unmarshaling JSON response: %v, %s", err, apiURL)
		return "", err
	}

	if len(mbResponse.Releases) > 0 {
		bestMatchingAlbum := getBestFromSearchResults(&mbResponse, nowPlaying)
		return bestMatchingAlbum, nil
	} else {
		return "", errors.New("No suitable album matches found while searching MusicBrainz")
	}

}

func foo2(nowPlaying *structs.NowPlaying) (string, error) {
	mbid, err := findAlbumMatch(nowPlaying)
	if err != nil {
		return "", err
	}

	artURL, err := getAlbumArtURL(mbid)
	if err != nil {
		return "", err
	}

	err = downloadAlbumArt(&artURL, &mbid)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/album-art/%s.jpg", mbid), nil
}

// func getAlbumArt(nowPlaying *structs.NowPlaying) {
// 	matchingMBID, err := findAlbumMatch(nowPlaying)
// 	if err != nil {

// 	}

// }
