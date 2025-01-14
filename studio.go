package scratchgonnect

import (
	"encoding/json"
	"io"
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

func (s Studio) GetCurators() UserArray {
	resp, err := http.Get("https://api.scratch.mit.edu/studios/" + to_string(s.Id) + "/curators")
	if err != nil {
		panic(err)
	}

	decoded := UserArray{}

	json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded
}

func (s Studio) GetManagers() UserArray {
	resp, err := http.Get("https://api.scratch.mit.edu/studios/" + to_string(s.Id) + "/managers")
	if err != nil {
		panic(err)
	}

	decoded := UserArray{}

	json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded
}

func (studio Studio) Follow(session Session) {
	resp, err := change_follow_request_studio(studio, session, "add")

	if err != nil || resp.StatusCode != 200 {
		panic("Follower action failed! http response:" + to_string(resp.StatusCode))
	}
}

func (studio Studio) Unfollow(session Session) {
	resp, err := change_follow_request_studio(studio, session, "remove")

	if err != nil || resp.StatusCode != 200 {
		panic("Follower action failed! http response:" + to_string(resp.StatusCode))
	}
}

// Functions

func change_follow_request_studio(s Studio, session Session, request string) (*http.Response, error) {
	req, err := http.NewRequest("PUT", "https://scratch.mit.edu/site-api/users/bookmarkers/"+to_string(s.Id)+"/"+request+"/?usernames="+session.Username, *new(io.Reader))
	if err != nil {
		panic(err)
	}

	req.Header = session.HttpHeader

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	req.AddCookie(&session.Cookie)
	req.AddCookie(&csrfCookieDefault)

	return http.DefaultClient.Do(req)
}

func GetStudio(studio string) *Studio {
	resp, err := http.Get("https://api.scratch.mit.edu/studios/" + studio)
	if err != nil {
		panic(err)
	}

	responseUser := new(Studio)

	json.NewDecoder(resp.Body).Decode(responseUser)

	return responseUser
}
