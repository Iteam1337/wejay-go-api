package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"sort"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	graphql "github.com/graph-gophers/graphql-go"
)

var mainSchema *graphql.Schema

// Schema is our GraphQL schema
var Schema = `
type Album {
  images: [Cover!]!
  name: String!
  uri: String!
}

type Artist {
  name: String!
  uri: String!
}

type Cover {
  height: Int!
  url: String!
  width: Int!
}

type Track {
  album: Album!
  artists: [Artist!]!
  duration: Int!
  name: String!
  spotifyUri: String!
  started: Int!
  user: User!
}

type User {
  email: String!
  id: String!
  lastPlay: Int!
}

type Room {
  name: String!
  currentTrack: Track
  queue: [Track]!
  users: [User]!
}

input QueueInput {
  roomName: String!
  spotifyId: [String!]!
  userId: String!
}

input JoinRoomInput {
  roomName: String!
  email: String!
}

type SearchResults {
  artists: [Artist]!
  tracks: [Track]!
}

type Query {
  rooms: [Room!]!
  room(name: String!): Room!
}

type Mutation {
  roomCreate(name: String!): Room
  roomJoin(input: JoinRoomInput!): Room
  roomNextTrack(roomName: String!): Room
  roomQueueTrack(input: QueueInput!): Room
  search(query: String!, limit: Int = 10): SearchResults!
}

# type Subscription {
#   onNextTrack(roomName: String!): Track
#   onPause(roomName: String!): Boolean
#   onPlay(roomName: String!): Boolean
#   roomUpdated(roomName: String!): Room
# }

schema {
  query: Query
  mutation: Mutation
}

`

var (
	// QueryNameNotProvided is thrown when a name is not provided
	QueryNameNotProvided = errors.New("no query was provided in the HTTP body")
)

// Handler fixes stuff
func Handler(context context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	// If no query is provided in the HTTP request body, throw an error
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, QueryNameNotProvided
	}

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	if err := json.Unmarshal([]byte(request.Body), &params); err != nil {
		log.Print("Could not decode body", err)
	}

	response := mainSchema.Exec(context, params.Query, params.OperationName, params.Variables)

	responseJSON, err := json.Marshal(response)

	if err != nil {
		log.Print("Could not decode body")
	}

	return events.APIGatewayProxyResponse{
		Body: string(responseJSON),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		StatusCode: 200,
	}, nil
}

func init() {
	mainSchema = graphql.MustParseSchema(Schema, &Resolver{})
}

func main() {
	lambda.Start(Handler)
}

func unique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func createQueue(tracks []*Track) []*Track {
	var users []string

	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].user.lastPlay < tracks[j].user.lastPlay
	})

	for _, track := range tracks {
		users = append(users, track.user.email)
	}

	uniqueUsers := unique(users)
	var sortedTracks [][]*Track
	var output []*Track

	for _, user := range uniqueUsers {
		var userTracks []*Track

		for _, track := range tracks {
			if track.user.email == user {
				userTracks = append(userTracks, track)
			}
		}

		sortedTracks = append(sortedTracks, userTracks)
	}

	maxLength := 0

	for _, user := range sortedTracks {
		sort.Slice(user, func(i, j int) bool {
			return user[i].added < user[j].added
		})

		if len(user) > maxLength {
			maxLength = len(user)
		}
	}

	for i := 0; i < maxLength; i++ {
		for _, tracks := range sortedTracks {
			if len(tracks) > i {
				output = append(output, tracks[i])
			}
		}
	}

	return output
}

func handleQueues() {
	for _, room := range rooms {
		if room.currentTrack != nil {
			currentTime := int32(time.Now().Unix()) * 1000
			trackEnd := room.currentTrack.started*1000 + room.currentTrack.duration

			if currentTime >= trackEnd {
				if len(room.queue) > 0 {
					room.currentTrack = room.queue[0]
					room.currentTrack.started = int32(time.Now().Unix())
					room.queue = room.queue[1:]
				} else {
					room.currentTrack = nil
				}
			}
		}
	}
}

func getSchema(path string) (string, error) {
	b, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
