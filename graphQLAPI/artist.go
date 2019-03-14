package main

// Artist struct
type Artist struct {
	name string
	uri  string
}

// Name resolves name field of Artist
func (a *Artist) Name() string {
	return a.name
}

// URI resolves uri field of Artist
func (a *Artist) URI() string {
	return a.uri
}
