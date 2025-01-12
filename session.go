package scratchgonnect

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Structs

type Session struct {
	Username   string
	Token      string
	Cookie     []string
	HttpHeader http.Header
}

type UserActivity struct {
	Id          int    `json:"id"`
	CreatedDate string `json:"datetime_created"`
	Username    string `json:"actor_username"`
	UserId      int    `json:"actor_id"`
	ProjectId   int    `json:"project_id"`
	ProjectName string `json:"project_title"`
	Type        string `json:"type"`
}

type UserActivityList []struct {
	UserActivity
}

// Struct functions

func (s Session) GetWhatsHappening() *UserActivityList {
	req, err := http.NewRequest("GET", "https://api.scratch.mit.edu/users/"+s.Username+"/following/users/activity", *new(io.Reader))
	if err != nil {
		panic(err)
	}

	req.Header = s.HttpHeader

	resp, _ := http.DefaultClient.Do(req)

	decoded := new(UserActivityList)
	json.NewDecoder(resp.Body).Decode(decoded)

	return decoded
}

func NewSession(username string, password string) *Session {
	post := JSONPostLogin{
		Username: username,
		Password: password,
	}

	b, err := json.Marshal(post)
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest("POST", "https://scratch.mit.edu/login/", bytes.NewReader(b))

	// Set header and cookie for login request

	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Add("user-agent", "(KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36")
	req.Header.Add("x-csrftoken", "a")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("referer", "https://scratch.mit.edu")
	req.Header.Add("Cookie", "scratchcsrftoken=a;scratchlanguage=en;")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	newToken := new(TokenDelivery)
	json.NewDecoder(resp.Body).Decode(newToken)

	// Create header and cookie for session requests

	defHeader := http.Header{}

	defHeader.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	defHeader.Add("user-agent", "(KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36")
	defHeader.Add("x-csrftoken", "a")
	defHeader.Add("x-requested-with", "XMLHttpRequest")
	defHeader.Add("referer", "https://scratch.mit.edu")
	defHeader.Add("X-Token", newToken[0].Token)

	defHeader.Add("Cookie", resp.Header["Set-Cookie"][0])
	defHeader.Add("Cookie", resp.Header["Set-Cookie"][1])

	// Create new session

	newSession := Session{
		Token:      newToken[0].Token,
		Username:   username,
		Cookie:     resp.Header["Set-Cookie"],
		HttpHeader: defHeader,
	}

	return &newSession
}
