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
		// handle error
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
			// handle error
		}

		req.Header.Set("Authorization", fmt.Sprintf("Token %s", listenbrainz_token))
		res, err := client.Do(req) // send the request
		if err != nil {
			// handle error
		}

		body, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		svc1 := structs.ListenbrainzResponse{}
		jsonErr := json.Unmarshal(body, &svc1)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		res1 := structs.NowPlaying{
			Artist: svc1.Payload.Listens[0].TrackMetadata.ArtistName,
			Album:  svc1.Payload.Listens[0].TrackMetadata.ReleaseName,
			Track:  svc1.Payload.Listens[0].TrackMetadata.TrackName,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res1)
	})

	http.ListenAndServe(":3000", r)
}
