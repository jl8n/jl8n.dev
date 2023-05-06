package main

import (
	"backend/structs"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var lbToken, lbUser string
var lastfmToken, lastfmUser string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables: ", err)
	}

	listenbrainz_token := os.Getenv("LISTENBRAINZ_API_TOKEN")
	listenbrainz_user := os.Getenv("LISTENBRAINZ_USER")
	lastfm_token := os.Getenv("LASTFM_API_TOKEN")
	lastfm_user := os.Getenv("LASTFM_USER")
	lbToken, lbUser = listenbrainz_token, listenbrainz_user
	lastfmToken, lastfmUser = lastfm_token, lastfm_user
}

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

func getAlbumMBID(w http.ResponseWriter, nowPlaying *structs.NowPlaying) {
	client := &http.Client{}
	// TODO: this looks bad
	apiURL := "http://musicbrainz.org/ws/2/release/?query=artist:" + url.QueryEscape(nowPlaying.Artist) + "release:" + url.QueryEscape(nowPlaying.Album) + "&fmt=json&limit=1"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("Error creating HTTP request: ", err)
		http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("User-Agent", "AlbumartCollage/0.1 ( josh.l8n@gmail.com )")
	res, err := client.Do(req) // send the request
	if err != nil {
		log.Println("Error sending HTTP request: ", err)
		http.Error(w, "Error sending `Now Playing` request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body2, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Println("Error reading HTTP response body: ", readErr)
		http.Error(w, "Error reading `Now Playing` response", http.StatusInternalServerError)
		return
	}

	mbResponse := structs.MBAlbumInfo{}
	jsonErr2 := json.Unmarshal(body2, &mbResponse)
	if jsonErr2 != nil {
		log.Println("Error unmarshaling JSON response: ", jsonErr2, apiURL)
		http.Error(w, "Error parsing `Now Playing 2` response", http.StatusInternalServerError)
		return
	}

	if len(mbResponse.Releases) > 0 {
		nowPlaying.AlbumMBID = mbResponse.Releases[0].ID // modify the reference
		fmt.Println(apiURL)
		fmt.Println(nowPlaying.AlbumMBID)
	} else {
		// playing song info not found in MusicBrainz
		w.WriteHeader(http.StatusNoContent)
	}
}

func getCoverArtURL(w http.ResponseWriter, nowPlaying *structs.NowPlaying) {
	client := &http.Client{} // TODO: should I include this in every function or pass it?
	apiURL := fmt.Sprintf("https://coverartarchive.org/release/%s", nowPlaying.AlbumMBID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("Error creating HTTP request: ", err)
		http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("User-Agent", "AlbumArtCollage/0.1 ( josh.l8n@gmail.com )")
	res, err := client.Do(req) // send the request
	if err != nil {
		log.Println("Error sending HTTP request: ", err)
		http.Error(w, "Error sending `Now Playing` request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("covert art not found, probably", res.StatusCode)
		return
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Println("Error reading HTTP response body: ", readErr)
		http.Error(w, "Error reading `Now Playing` response", http.StatusInternalServerError)
		return
	}

	caResponse := structs.CoverartArchive{}
	jsonErr := json.Unmarshal(body, &caResponse)
	if jsonErr != nil {
		log.Println("Error unmarshaling JSON response: ", jsonErr, apiURL)
		http.Error(w, "Error parsing `Now Playing` response", http.StatusInternalServerError)
		return
	}

	// TODO: fallback to smaller sizes if large not found
	nowPlaying.Coverart = caResponse.Images[0].Image
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		//AllowOriginFunc:  AllowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		res1 := structs.NowPlaying{}
		fmt.Println("GET Artist, Album, Track")
		getNowPlaying(w, &res1) // add artist, album, and track to NowPlaying{}
		fmt.Println("GET Album MBID")
		getAlbumMBID(w, &res1) // add album mbid to NowPlaying{}
		fmt.Println("GET Coverart")
		getCoverArtURL(w, &res1)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res1)
	})

	// r.Get("/nowplaying", func(w http.ResponseWriter, r *http.Request) {

	// 	client := &http.Client{}
	// 	url := fmt.Sprintf("https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=%s&limit=1&extended=1&nowplaying=true&api_key=%s&format=json", lastfmUser, lastfmToken)

	// 	req, err := http.NewRequest("GET", url, nil)
	// 	if err != nil {
	// 		log.Println("Error creating HTTP request: ", err)
	// 		http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	res, err := client.Do(req) // send the request
	// 	if err != nil {
	// 		log.Println("Error sending HTTP request: ", err)
	// 		http.Error(w, "Error sending `Now Playing` request", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	body, readErr := io.ReadAll(res.Body)
	// 	if readErr != nil {
	// 		log.Println("Error reading HTTP response body: ", readErr)
	// 		http.Error(w, "Error reading `Now Playing` response", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	svc1 := structs.LastfmNowPlaying{}
	// 	jsonErr := json.Unmarshal(body, &svc1)
	// 	if jsonErr != nil {
	// 		log.Println("Error unmarshaling JSON response: ", jsonErr)
	// 		http.Error(w, "Error parsing `Now Playing` response", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	if len(svc1.Recenttracks.Track) > 0 {
	// 		// loop over recent tracks to find if one is currently playing
	// 		fmt.Println("it here")
	// 		for _, track := range svc1.Recenttracks.Track {
	// 			fmt.Println(track)
	// 			if track.Attr.Nowplaying == "true" {
	// 				res1 := structs.NowPlaying{
	// 					Artist:    track.Artist.Name,
	// 					Album:     track.Album.Text,
	// 					Track:     track.Name,
	// 					AlbumMBID: track.Album.Mbid,
	// 				}
	// 				w.Header().Set("Content-Type", "application/json")
	// 				w.WriteHeader(http.StatusCreated)
	// 				json.NewEncoder(w).Encode(res1)
	// 				return
	// 			}
	// 		}

	// 		w.WriteHeader(http.StatusNoContent)
	// 		return

	// 	} else {
	// 		http.Error(w, "No data returned by API", http.StatusInternalServerError)
	// 		return
	// 	}
	// })

	http.ListenAndServe(":3000", r)
}
