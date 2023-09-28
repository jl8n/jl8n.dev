package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

	r.Get("/nowplaying", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		nowPlaying, err := getCurrentlyPlayingSong()
		if err != nil {
			fmt.Fprintf(w, "data: %s\n\n", "Error retrieving currently playing song")
			flusher.Flush()
			return
		}

		fmt.Println(nowPlaying.Track, nowPlaying.Album, nowPlaying.Artist)

		if nowPlaying.Track == "" {
			fmt.Fprintf(w, "data: %s\n\n", "No song is currently playing")
			flusher.Flush()
			return
		}

		// send server-side events to the client
		fmt.Fprintf(w, "event: Track\ndata: %s\n\n", nowPlaying.Track)
		fmt.Fprintf(w, "event: Album\ndata: %s\n\n", nowPlaying.Album)
		fmt.Fprintf(w, "event: Artist\ndata: %s\n\n", nowPlaying.Artist)

		flusher.Flush()

		time.Sleep(time.Second * 10)

		// albumArt := getAlbumArt(nowPlaying.Album)
		// if albumArt != "" {
		// 	fmt.Fprintf(w, "data: %s\n\n", albumArt)
		// 	flusher.Flush()
		// }
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		nowPlaying, err := getNowPlaying()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
			return
		}

		if nowPlaying.Track == "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// create context to ensure that isReleaseSavedInDB() doesn't time out
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if isReleaseSavedInDB(ctx, &nowPlaying) {
			fmt.Println("release exists", nowPlaying)
			nowPlaying.ArtAvailable = true
			sendResponse(w, nowPlaying)
			return
		}

		fmt.Println("release does not exist", nowPlaying)

		foundMBID, err := getAlbumMBID(w, &nowPlaying) // add album mbid to NowPlaying{}
		if err != nil {
			fmt.Println(err)
			sendResponse(w, nowPlaying)
			return
		}

		nowPlaying.AlbumMBID = foundMBID

		artURL, err := getCoverArtURL(w, &nowPlaying)
		fmt.Println("art url", artURL)
		if err != nil {
			// cover art not found via search - respond with metadata only
			fmt.Println("Could not get album art URL from coverartarchive.org")
			sendResponse(w, nowPlaying)
			return
		}

		err = downloadAlbumArt(artURL)
		if err != nil {
			// album art not downloaded - respond with metadata only
			//  TODO: respond with status indicating that album art is unavailable
			sendResponse(w, nowPlaying)
			return
		}

		err = addToMongo(&nowPlaying)
		if err != nil {
			// could not add release to MongoDB
			// only effect of this is slower subsequent
		}
		nowPlaying.ArtAvailable = true

		sendResponse(w, nowPlaying)
	})

	r.Get("/albumart", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET /albumart")
		filename := r.URL.Query().Get("filename")
		if filename == "" {
			http.Error(w, "filename parameter is required", http.StatusBadRequest)
			return
		}

		filepath := "album-art/" + strings.TrimSpace(filename) + ".jpg"
		fmt.Println(filepath)
		http.ServeFile(w, r, filepath)
	})

	r.Get("/events", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET /events")
		// Set the necessary headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Create a channel to send messages to the client
		messageChan := make(chan string)

		go func() {
			for {
				time.Sleep(time.Second * 2)
				messageChan <- time.Now().Format(time.RFC1123)
			}
		}()

		for {
			select {
			case msg := <-messageChan:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			case <-r.Context().Done():
				return
			}
		}
	})

	http.ListenAndServe(":3000", r)
}
