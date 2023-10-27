package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var lbToken, lbUser string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables: ", err)
	}

	listenbrainz_token := os.Getenv("LISTENBRAINZ_API_TOKEN")
	listenbrainz_user := os.Getenv("LISTENBRAINZ_USER")
	lbToken, lbUser = listenbrainz_token, listenbrainz_user
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		//AllowOriginFunc:  AllowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "AuthFlusherorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Serve static files from the "album-art" directory
	fs := http.StripPrefix("/album-art", http.FileServer(http.Dir("../album-art")))
	r.Handle("/album-art/*", fs)

	r.Get("/nowplaying", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Server-sent events are unsupported!", http.StatusInternalServerError)
			return
		}

		nowPlaying, err := getCurrentlyPlayingSong()
		if err != nil {
			fmt.Fprintf(w, "data: %s\n\n", "Error retrieving currently playing song")
			flusher.Flush()
			return
		}

		setupSSEHeaders(w)
		getSongAndSendSSE(w, flusher, &nowPlaying)

		// create a channel to receive the album-art URL from the goroutine
		albumArtChan := make(chan string)

		// create a context to control goroutine timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// search for and download album art (in a separate thread)
		go func() {
			select {
			case <-ctx.Done():
				return // The context was cancelled or timed out
			default:
				path, err := downloadAlbumArt(&nowPlaying)
				if err != nil {
					albumArtChan <- fmt.Sprintf("Error downloading album art: %s", err.Error())
					return
				}

				albumArtChan <- path
			}
		}()

		// either wait for the album art or timeout after 10 seconds
		select {
		case albumArt := <-albumArtChan:
			if strings.HasPrefix(albumArt, "Error") {
				fmt.Fprintf(w, "data: %s\n\n", albumArt) // send SSE with error message
			} else if albumArt != "" {
				fmt.Fprintf(w, "event: AlbumArt\ndata: %s\n\n", albumArt) // send SSE with album art
			}
			flusher.Flush()
		case <-time.After(time.Second * 10):
			// no album art was found within 10 seconds
		}
	})

	type Response struct {
		Files []string `json:"files"`
	}

	r.Get("/art", func(w http.ResponseWriter, r *http.Request) {
		var res Response

		files, _ := os.ReadDir("../album-art")
		for _, f := range files {
			res.Files = append(res.Files, "http://localhost:3000/album-art/"+f.Name())
		}

		json.NewEncoder(w).Encode(res)
	})

	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	nowPlaying, err := getNowPlaying()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	if nowPlaying.Track == "" {
	// 		w.WriteHeader(http.StatusNoContent)
	// 		return
	// 	}

	// 	// create context to ensure that isReleaseSavedInDB() doesn't time out
	// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 	defer cancel()

	// 	if isReleaseSavedInDB(ctx, &nowPlaying) {
	// 		fmt.Println("release exists", nowPlaying)
	// 		nowPlaying.ArtAvailable = true
	// 		sendResponse(w, nowPlaying)
	// 		return
	// 	}

	// 	fmt.Println("release does not exist", nowPlaying)

	// 	foundMBID, err := getAlbumMBID(w, &nowPlaying) // add album mbid to NowPlaying{}
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		sendResponse(w, nowPlaying)
	// 		return
	// 	}

	// 	nowPlaying.AlbumMBID = foundMBID

	// 	artURL, err := getCoverArtURL(w, &nowPlaying)
	// 	fmt.Println("art url", artURL)
	// 	if err != nil {
	// 		// cover art not found via search - respond with metadata only
	// 		fmt.Println("Could not get album art URL from coverartarchive.org")
	// 		sendResponse(w, nowPlaying)
	// 		return
	// 	}

	// 	err = downloadAlbumArt1(artURL)
	// 	if err != nil {
	// 		// album art not downloaded - respond with metadata only
	// 		//  TODO: respond with status indicating that album art is unavailable
	// 		sendResponse(w, nowPlaying)
	// 		return
	// 	}

	// 	err = addToMongo(&nowPlaying)
	// 	if err != nil {
	// 		// could not add release to MongoDB
	// 		// only effect of this is slower subsequent
	// 	}
	// 	nowPlaying.ArtAvailable = true

	// 	sendResponse(w, nowPlaying)
	// })

	// r.Get("/albumart", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("GET /albumart")
	// 	filename := r.URL.Query().Get("filename") // url query parameter
	// 	if filename == "" {
	// 		http.Error(w, "filename parameter is required", http.StatusBadRequest)
	// 		return
	// 	}

	// 	filepath := "album-art/" + strings.TrimSpace(filename) + ".jpg"
	// 	fmt.Println(filepath)
	// 	http.ServeFile(w, r, filepath)
	// })

	// r.Get("/events", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("GET /events")
	// 	// Set the necessary headers for SSE
	// 	w.Header().Set("Content-Type", "text/event-stream")
	// 	w.Header().Set("Cache-Control", "no-cache")
	// 	w.Header().Set("Connection", "keep-alive")

	// 	// Create a channel to send messages to the client
	// 	messageChan := make(chan string)

	// 	go func() {
	// 		for {
	// 			time.Sleep(time.Second * 2)
	// 			messageChan <- time.Now().Format(time.RFC1123)
	// 		}
	// 	}()

	// 	for {
	// 		select {
	// 		case msg := <-messageChan:
	// 			fmt.Fprintf(w, "data: %s\n\n", msg)
	// 			if f, ok := w.(http.Flusher); ok {
	// 				f.Flush()
	// 			}
	// 		case <-r.Context().Done():
	// 			return
	// 		}
	// 	}
	// })

	port := 3000
	fmt.Printf("Starting server on port %d...\n", port)

	http.ListenAndServe(":3000", r)
}
