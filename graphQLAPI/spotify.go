package main

import (
	"context"
	"log"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

var client spotify.Client

func init() {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())

	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client = spotify.Authenticator{}.NewClient(token)
}

// SpotifyGetTrack gets a track from the Spotify's web API
func SpotifyGetTrack(trackID spotify.ID) *spotify.FullTrack {
	results, err := client.GetTrack(trackID)

	if err != nil {
		log.Fatal(err)
	}

	return results
}

// SpotifySearchTrack searches after tracks in the Spotify web API
func SpotifySearchTrack(query string, limit int32) *spotify.SearchResult {
	results, err := client.SearchOpt(query, spotify.SearchTypeTrack|spotify.SearchTypeArtist, &spotify.Options{Limit: createInt(limit)})

	if err != nil {
		log.Fatal(err)
	}

	return results
}

func createInt(x int32) *int {
	val := int(x)
	return &val
}
