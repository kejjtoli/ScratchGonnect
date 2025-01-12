package scratchgonnect

import (
	"encoding/json"
	"net/http"
)

// Struct declaration

type Studio struct {
	Id              int    `json:"id"`
	Name            string `json:"title"`
	HostId          int    `json:"host"`
	Description     string `json:"description"`
	Visiblity       string `json:"visibility"`
	IsPublic        bool   `json:"public"`
	IsOpen          bool   `json:"open_to_all"`
	CommentsAllowed bool   `json:"comments_allowed"`
	Image           string `json:"image"`
	History         struct {
		JoinedDate   string `json:"joined"`
		ModifiedDate string `json:"modified"`
	} `json:"history"`
	Stats struct {
		Comments  int `json:"comments"`
		Followers int `json:"followers"`
		Managers  int `json:"managers"`
		Projects  int `json:"projects"`
	} `json:"stats"`
}

// Struct functions

func (s Studio) GetCurators() *UserArray {
	resp, err := http.Get("https://api.scratch.mit.edu/studios/" + to_string(s.Id) + "/curators")
	if err != nil {
		panic(err)
	}

	decoded := new(UserArray)

	json.NewDecoder(resp.Body).Decode(decoded)

	return decoded
}

func (s Studio) GetManagers() *UserArray {
	resp, err := http.Get("https://api.scratch.mit.edu/studios/" + to_string(s.Id) + "/managers")
	if err != nil {
		panic(err)
	}

	decoded := new(UserArray)

	json.NewDecoder(resp.Body).Decode(decoded)

	return decoded
}

// Functions

func GetStudio(studio string) *Studio {
	resp, err := http.Get("https://api.scratch.mit.edu/studios/" + studio)
	if err != nil {
		panic(err)
	}

	responseUser := new(Studio)

	json.NewDecoder(resp.Body).Decode(responseUser)

	return responseUser
}
