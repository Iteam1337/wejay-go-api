package main

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/graph-gophers/graphql-go/errors"
	"github.com/zmb3/spotify"
)

var rooms []*Room

// Resolver is the root resolver
type Resolver struct{}

// Room resolves room query
func (r *Resolver) Room(args struct{ Name string }) (*Room, error) {
	room := findRoom(args.Name)

	if room == nil {
		return nil, &errors.QueryError{Message: "Room does not exist"}
	}

	return room, nil
}

// Rooms resolves rooms query
func (r *Resolver) Rooms() []*Room {
	return rooms
}

// RoomCreate resolves addRoom mutation
func (r *Resolver) RoomCreate(args struct{ Name string }) (*Room, error) {
	currentRoom := findRoom(args.Name)

	if currentRoom != nil {
		return nil, &errors.QueryError{Message: "Room already exists"}
	}

	newRoom := &Room{name: args.Name}

	rooms = append(rooms, newRoom)

	return newRoom, nil
}

// RoomNextTrack resolves roomNextRoom mutation
func (r *Resolver) RoomNextTrack(args struct{ RoomName string }) (*Room, error) {
	room := findRoom(args.RoomName)

	if room == nil {
		return nil, &errors.QueryError{Message: "Room does not exist"}
	}

	if len(room.queue) == 0 {
		room.currentTrack = nil

		return room, nil
	}

	room.currentTrack = room.queue[0]
	room.currentTrack.started = int32(time.Now().Unix())
	room.queue = room.queue[1:]

	return room, nil
}

type queueInput struct {
	RoomName  string
	SpotifyID []spotify.ID
	UserID    string
}

// RoomQueueTrack resolves queueTrack mutation
func (r *Resolver) RoomQueueTrack(args struct{ Input queueInput }) (*Room, error) {
	room := findRoom(args.Input.RoomName)

	for _, trackID := range args.Input.SpotifyID {
		trackIDAsURI := "spotify:track:" + string(trackID)

		if room == nil {
			return nil, &errors.QueryError{Message: "Room does not exist"}
		}

		for _, track := range room.queue {
			if string(track.spotifyURI) == trackIDAsURI {
				return nil, &errors.QueryError{Message: "Already in queue"}
			}
		}

		if room.currentTrack != nil && (string(room.currentTrack.spotifyURI) == trackIDAsURI) {
			return nil, &errors.QueryError{Message: "Currently playing track"}
		}

		currentUser := findUser(args.Input.UserID, room)

		if currentUser == nil {
			return nil, &errors.QueryError{Message: "User is not in room"}
		}

		results := SpotifyGetTrack(trackID)

		newTrack := &Track{
			added:      time.Now().Unix(),
			album:      results.Album,
			artists:    results.Artists,
			duration:   int32(results.Duration),
			name:       results.Name,
			spotifyURI: results.URI,
			user:       currentUser,
		}

		currentUser.lastPlay = time.Now().Unix()
		room.queue = append(room.queue, newTrack)
	}

	room.queue = createQueue(room.queue)

	if room.currentTrack == nil {
		room.currentTrack = room.queue[0]
		room.currentTrack.started = int32(time.Now().Unix())
		room.queue = room.queue[1:]
	}

	return room, nil
}

type joinInput struct {
	Email    string
	RoomName string
}

// RoomJoin handles joining a room
func (r *Resolver) RoomJoin(args struct{ Input joinInput }) (*Room, error) {
	currentRoom := findRoom(args.Input.RoomName)

	if currentRoom == nil {
		return nil, &errors.QueryError{Message: "Room does not exist"}
	}

	for _, user := range currentRoom.users {
		if user.email == args.Input.Email {
			return nil, &errors.QueryError{Message: "Already in the room"}
		}
	}

	newUser := &User{
		email:    args.Input.Email,
		id:       GetMD5Hash(args.Input.Email),
		lastPlay: 0.0,
	}

	currentRoom.users = append(currentRoom.users, newUser)

	return currentRoom, nil
}

type searchArgs struct {
	Query string
	Limit int32
}

// Search resolves search mutation
func (r *Resolver) Search(args searchArgs) (*SearchResults, error) {
	results := SpotifySearchTrack(args.Query, args.Limit)
	var trackResults []*Track
	var artistResults []*Artist

	if results.Tracks != nil {
		for _, track := range results.Tracks.Tracks {
			newTrack := &Track{
				album:      track.Album,
				artists:    track.Artists,
				duration:   int32(track.Duration),
				name:       track.Name,
				spotifyURI: track.URI,
			}

			trackResults = append(trackResults, newTrack)
		}
	}

	if results.Artists != nil {
		for _, artist := range results.Artists.Artists {
			newArtist := &Artist{
				name: artist.Name,
				uri:  string(artist.URI),
			}

			artistResults = append(artistResults, newArtist)
		}
	}

	return &SearchResults{artists: artistResults, tracks: trackResults}, nil
}

func roomNextTrack(roomName string) (*Room, error) {
	room := findRoom(roomName)

	if room == nil {
		return nil, &errors.QueryError{Message: "Room does not exist"}
	}

	if len(room.queue) == 0 {
		room.currentTrack = nil

		return room, nil
	}

	room.currentTrack = room.queue[0]
	room.currentTrack.started = int32(time.Now().Unix())
	room.queue = room.queue[1:]

	return room, nil
}

// GetMD5Hash creates an md5 hash from a string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func findRoom(name string) *Room {
	for _, room := range rooms {
		if room.name == name {
			return room
		}
	}

	return nil
}

func findUser(userID string, room *Room) *User {
	for _, user := range room.users {
		if user.id == userID {
			return user
		}
	}

	return nil
}
