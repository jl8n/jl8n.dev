package main

import (
	"backend/structs"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var listenbrainz_token, listenbrainz_user string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables: ", err)
	}

	token := os.Getenv("LISTENBRAINZ_API_TOKEN")
	user := os.Getenv("LISTENBRAINZ_USER")
	listenbrainz_token, listenbrainz_user = token, user
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

		client := &http.Client{}
		url := fmt.Sprintf("https://api.listenbrainz.org/1/user/%s/playing-now", listenbrainz_user)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("Error creating HTTP request: ", err)
			http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Authorization", fmt.Sprintf("Token %s", listenbrainz_token))
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

		svc1 := structs.ListenbrainzResponse{}
		jsonErr := json.Unmarshal(body, &svc1)
		if jsonErr != nil {
			log.Println("Error unmarshaling JSON response: ", jsonErr)
			http.Error(w, "Error parsing `Now Playing` response", http.StatusInternalServerError)
			return
		}

		if len(svc1.Payload.Listens) > 0 {
			res1 := structs.NowPlaying{
				Artist: svc1.Payload.Listens[0].TrackMetadata.ArtistName,
				Album:  svc1.Payload.Listens[0].TrackMetadata.ReleaseName,
				Track:  svc1.Payload.Listens[0].TrackMetadata.TrackName,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(res1)
		} else {
			http.Error(w, "No data returned by API", http.StatusInternalServerError)
		}
	})

	r.Get("/listens", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}
		url := fmt.Sprintf("https://api.listenbrainz.org/1/user/%s/playing-now", listenbrainz_user)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("Error creating HTTP request: ", err)
			http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Authorization", fmt.Sprintf("Token %s", listenbrainz_token))
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
	})
	http.ListenAndServe(":3000", r)
}
