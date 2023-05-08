package main

import (
	"backend/structs"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

		// create context to ensure that isReleaseExists() doesn't time out
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if isReleaseExists(ctx, &res1) {
			fmt.Println("release exists", res1)
			res1.ArtAvailable = true
			sendResponse(w, res1)
			return
		}

		fmt.Println("release does not exist", res1)

		foundMBID, err := getAlbumMBID(w, &res1) // add album mbid to NowPlaying{}
		if err != nil {
			fmt.Println(err)
			sendResponse(w, res1)
			return
		}

		res1.AlbumMBID = foundMBID

		artURL, err := getCoverArtURL(w, &res1)
		fmt.Println("art url", artURL)
		if err != nil {
			// cover art not found via search - respond with metadata only
			fmt.Println("Could not get album art URL from coverartarchive.org")
			sendResponse(w, res1)
			return
		}

		err = downloadAlbumArt(artURL)
		if err != nil {
			// album art not downloaded - respond with metadata only
			//  TODO: respond with status indicating that album art is unavailable
			sendResponse(w, res1)
			return
		}

		err = addToMongo(&res1)
		if err != nil {
			// could not add release to MongoDB
			// only effect of this is slower subsequent
		}
		res1.ArtAvailable = true

		sendResponse(w, res1)
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

	http.ListenAndServe(":3000", r)
}

func sendResponse(w http.ResponseWriter, res structs.NowPlaying) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func isReleaseExists(ctx context.Context, nowPlaying *structs.NowPlaying) bool {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://root:example@localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		// handle error
	}
	defer client.Disconnect(ctx)

	// Check the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		// handle error
	}

	// Get a handle to the releases collection
	collection := client.Database("scrobbles").Collection("releases")

	filter := bson.M{"name": nowPlaying.Album, "artist": nowPlaying.Artist}
	var result bson.M // create variable to store result in
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Document doesn't exist
			return false
		} else {
			// Error occurred TODO: use `err` var for error handling
			return false
		}
	}

	// Document exists, return its _id
	id := result["_id"].(string)
	nowPlaying.AlbumMBID = id
	return true
}

func addToMongo(nowPlaying *structs.NowPlaying) error {
	fmt.Println("Storing release...")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		return err
	}

	// Attempt to establish a connection or timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Release resources associated with the context

	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	coll := client.Database("scrobbles").Collection("releases")
	doc := structs.Release{
		ID:     nowPlaying.AlbumMBID,
		Name:   nowPlaying.Album,
		Artist: nowPlaying.Artist,
	}

	// Insert document into collection
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return nil
}

func getFilenameFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	segments := strings.Split(parsedURL.Path, "/")
	filename := segments[len(segments)-2] + ".jpg"
	return filename, nil
}

func downloadAlbumArt(artURL string) error {
	filename, err := getFilenameFromURL(artURL)
	if err != nil {
		return errors.New("Could not use MBID as filename")
	}

	// Create the directory if it doesn't exist
	err = os.MkdirAll("album-art", 0755)
	if err != nil {
		return errors.New("Could not create directory")
	}

	// Create the empty file on the filesystem
	out, err := os.Create(filepath.Join("album-art", filename))
	if err != nil {
		return errors.New("Could not create empty file on filesystem")
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(artURL)
	if err != nil {
		return errors.New("Could not download album art data")
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.New("Could not write album art data to file")
	}

	return nil
}
