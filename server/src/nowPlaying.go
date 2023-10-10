package main

import (
	"backend/structs"
	"fmt"
	"net/http"
)

// TODO:
// These functions are used by the /now-playing route (hence the filename)
// This file should probabaly be renamed and organized differently

func setupHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

func getSongAndSendSSE(w http.ResponseWriter, flusher http.Flusher, nowPlaying *structs.NowPlaying) {
	fmt.Printf("%s - %s - %s\n", nowPlaying.Artist, nowPlaying.Track, nowPlaying.Album)

	if nowPlaying.Track == "" {
		fmt.Fprintf(w, "data: %s\n\n", "No song is currently playing")
		flusher.Flush()
		return
	}

	// send server-sent events to the client
	fmt.Fprintf(w, "event: Track\ndata: %s\n\n", nowPlaying.Track)
	fmt.Fprintf(w, "event: Album\ndata: %s\n\n", nowPlaying.Album)
	fmt.Fprintf(w, "event: Artist\ndata: %s\n\n", nowPlaying.Artist)
	flusher.Flush()
}
