package main

// SearchResults takes care of search results
type SearchResults struct {
	artists []*Artist
	tracks  []*Track
}

// Artists resolves artists on SearchResults
func (s *SearchResults) Artists() []*Artist {
	return s.artists
}

// Tracks resolves tracks on SearchResults
func (s *SearchResults) Tracks() []*Track {
	return s.tracks
}
