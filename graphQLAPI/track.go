package main

import "github.com/zmb3/spotify"

// Track struct
type Track struct {
	added      int64
	album      spotify.SimpleAlbum
	artists    []spotify.SimpleArtist
	duration   int32
	name       string
	spotifyURI spotify.URI
	started    int32
	user       *User
}

// Album resolves album field of Track
func (t *Track) Album() *Album {
	return &Album{
		images: t.album.Images,
		name:   t.album.Name,
		uri:    string(t.album.URI),
	}
}

// Artists resolves artist field of Track
func (t *Track) Artists() []*Artist {
	var artists []*Artist

	for _, artist := range t.artists {
		artists = append(artists, &Artist{
			name: artist.Name,
			uri:  string(artist.URI),
		})
	}

	return artists
}

// Name resolves name field of Track
func (t *Track) Name() string {
	return t.name
}

// Duration resolves duration field of Track
func (t *Track) Duration() int32 {
	return t.duration
}

// SpotifyURI resolves spotify uri field of Track
func (t *Track) SpotifyURI() string {
	return string(t.spotifyURI)
}

// Started resolves started field of Track
func (t *Track) Started() int32 {
	return t.started
}

// User resolves user field of Track
func (t *Track) User() *User {
	return t.user
}
