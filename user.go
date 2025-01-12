package scratchgonnect

import (
	"encoding/json"
	"net/http"
)

// Structs

type User struct {
	Id            int    `json:"id"`
	Username      string `json:"username"`
	IsScratchteam bool   `json:"scratchteam"`
	History       struct {
		JoinedDate string `json:"joined"`
	} `json:"history"`
	Profile struct {
		Id      int    `json:"id"`
		Status  string `json:"status"`
		Bio     string `json:"bio"`
		Country string `json:"country"`
		Images  struct {
			Res90x90 string `json:"90x90"`
			Res60x60 string `json:"60x60"`
			Res55x55 string `json:"55x55"`
			Res50x50 string `json:"50x50"`
			Res32x32 string `json:"32x32"`
		}
	}
}

// Struct functions

func (u User) GetFollowers() *UserArray {
	resp, err := http.Get("https://api.scratch.mit.edu/users/" + u.Username + "/followers")
	if err != nil {
		panic(err)
	}

	decoded := new(UserArray)

	json.NewDecoder(resp.Body).Decode(decoded)

	return decoded
}

func (u User) GetFollowing() *UserArray {
	resp, err := http.Get("https://api.scratch.mit.edu/users/" + u.Username + "/following")
	if err != nil {
		panic(err)
	}

	decoded := new(UserArray)

	json.NewDecoder(resp.Body).Decode(decoded)

	return decoded
}

// Functions

func GetUser(username string) *User {
	resp, err := http.Get("https://api.scratch.mit.edu/users/" + username)
	if err != nil {
		panic(err)
	}

	responseUser := new(User)

	json.NewDecoder(resp.Body).Decode(responseUser)

	return responseUser
}
