package main

// Room struct
type Room struct {
	currentTrack *Track
	name         string
	queue        []*Track
	users        []*User
}

// CurrentTrack resolves currentTrack field of Room
func (r *Room) CurrentTrack() *Track {
	return r.currentTrack
}

// Name resolves name field of Room
func (r *Room) Name() string {
	return r.name
}

// Queue resolves queue field of Room
func (r *Room) Queue() []*Track {
	return r.queue
}

// Users resolves users field of Room
func (r *Room) Users() []*User {
	return r.users
}
