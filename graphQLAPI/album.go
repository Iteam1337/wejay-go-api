package main

import "github.com/zmb3/spotify"

// Album struct
type Album struct {
	images []spotify.Image
	name   string
	uri    string
}

// NAME resolves name field of Album
func (a *Album) NAME() string {
	return a.name
}

// IMAGES resolves images field of Album
func (a *Album) IMAGES() []*Cover {
	var covers []*Cover

	for _, cover := range a.images {
		covers = append(covers, &Cover{
			height: cover.Height,
			url:    cover.URL,
			width:  cover.Width,
		})
	}

	return covers
}

// URI resolves uri field of Album
func (a *Album) URI() string {
	return string(a.uri)
}
