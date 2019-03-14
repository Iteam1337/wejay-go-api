package main

// User is a user
type User struct {
	email    string
	id       string
	lastPlay int64
}

// Email resolves email field of User
func (u *User) Email() string {
	return u.email
}

// ID resolves id field of User
func (u *User) ID() string {
	return u.id
}

// LastPlay resolves lastPlay field of User
func (u *User) LastPlay() int32 {
	return int32(u.lastPlay)
}
