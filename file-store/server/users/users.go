package users

import (
	"errors"
	"math/rand"
	"strconv"
)

type User struct {
	ID       string
	Username string
	// TODO password salt and pepper. That's not the focus of the example.
	Password string
	// Set when sending the user to authorize, removed after we get the token.
	StorageAuthorizationCodePkce string
	// Set if we have it
	StorageToken string
}

var users = make(map[string]*User) // userId -> User

func CreateUser(username string, password string) (*User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password required")
	}
	for _, user := range users {
		if user.Username == username {
			return nil, errors.New("username already taken")
		}
	}
	user := User{strconv.FormatInt(rand.Int63(), 10), username, password, "", ""}
	users[user.ID] = &user
	return &user, nil
}

func LoginUser(username string, password string) (*User, error) {
	for _, user := range users {
		if username == user.Username && password == user.Password {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func GetUser(id string) (*User, error) {
	user, ok := users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *User) SetStorageAuthorizationPkce(pkce string) {
	u.StorageAuthorizationCodePkce = pkce
}

func (u *User) SetStorageToken(token string) {
	u.StorageAuthorizationCodePkce = ""
	u.StorageToken = token
}
