package main

import (
	"backend/structs"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func sendResponse(w http.ResponseWriter, res structs.NowPlaying) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func isReleaseSavedInDB(ctx context.Context, nowPlaying *structs.NowPlaying) bool {
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
	fmt.Println("\nattempting to connect")
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
	fmt.Println("\ndid not insert")
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

// func downloadAlbumArt1(artURL string) error {
// 	filename, err := getFilenameFromURL(artURL)
// 	if err != nil {
// 		return errors.New("could not use MBID as filename")
// 	}

// 	// Create the directory if it doesn't exist
// 	err = os.MkdirAll("album-art", 0755)
// 	if err != nil {
// 		return errors.New("could not create directory")
// 	}

// 	// Create the empty file on the filesystem
// 	out, err := os.Create(filepath.Join("album-art", filename))
// 	if err != nil {
// 		return errors.New("could not create empty file on filesystem")
// 	}
// 	defer out.Close()

// 	// Get the data
// 	resp, err := http.Get(artURL)
// 	if err != nil {
// 		return errors.New("could not download album art data")
// 	}
// 	defer resp.Body.Close()

// 	// Write the body to file
// 	_, err = io.Copy(out, resp.Body)
// 	if err != nil {
// 		return errors.New("could not write album art data to file")
// 	}

// 	return nil
// }
