package scratchgonnect

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
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

func (u User) GetFollowers() []User {
	resp, err := http.Get("https://api.scratch.mit.edu/users/" + u.Username + "/followers")
	if err != nil {
		panic(err)
	}

	decoded := []User{}

	json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded
}

func (u User) GetFollowing() []User {
	resp, err := http.Get("https://api.scratch.mit.edu/users/" + u.Username + "/following")
	if err != nil {
		panic(err)
	}

	decoded := []User{}

	json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded
}

func (u User) GetProjects() []Project {
	resp, err := http.Get("https://api.scratch.mit.edu/users/" + u.Username + "/projects")
	if err != nil || resp.StatusCode != 200 {
		panic("Get projects action failed! http response:" + to_string(resp.StatusCode))
	}

	decoded := []Project{}

	json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded
}

func (u User) GetUnreadMessageCount() int {
	resp, err := http.Get("https://api.scratch.mit.edu/users/" + u.Username + "/messages/count")
	if err != nil || resp.StatusCode != 200 {
		panic("Get projects action failed! http response:" + to_string(resp.StatusCode))
	}

	decoded := json_msgs{}

	json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded.Count
}

func (u User) Follow(session Session) {
	resp, err := change_follow_request(u, session, "add")

	if err != nil || resp.StatusCode != 200 {
		panic("Follower action failed! http response:" + to_string(resp.StatusCode))
	}
}

func (u User) Unfollow(session Session) {
	resp, err := change_follow_request(u, session, "remove")

	if err != nil || resp.StatusCode != 200 {
		panic("Follower action failed! http response:" + to_string(resp.StatusCode))
	}
}

func (u User) PostComment(session Session, content string, parent_id string, commentee_id string) {
	b, _ := json.Marshal(json_comment{
		Content:  content,
		ParentId: parent_id,
		Id:       commentee_id,
	})

	req, err := http.NewRequest("POST", "https://scratch.mit.edu/site-api/comments/user/"+u.Username+"/add/", bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}

	req.Header = session.HttpHeader

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	req.AddCookie(&session.Cookie)
	req.AddCookie(&csrfCookieDefault)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		panic("Post comment failed! http response:" + to_string(resp.StatusCode))
	}
}

func (u User) GetComments(page int) []Comment {
	resp, err := http.Get("https://scratch.mit.edu/site-api/comments/user/" + u.Username + "/?page=" + to_string(page))
	if err != nil {
		panic(err)
	}

	decoded := []Comment{}

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)
	comments := strings.Split(body, `<li class="top-level-reply">`)

	for _, c_string := range comments {
		elements := strings.Split(c_string, `<`)
		id := scrape_element(elements[1], "data-comment-id")

		if id != "" {
			username := scrape_element(elements[4], "data-comment-user")
			content := strings.Split(elements[12], `>`)[1]

			int_id, err := strconv.Atoi(id)
			if err != nil {
				panic(err)
			}

			new_comment := Comment{
				Id:      int_id,
				Content: strings.TrimSpace(content),
				Author: User{
					Username: username,
				},
			}

			decoded = append(decoded, new_comment)
		}
	}

	json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded
}

// Functions

func change_follow_request(u User, session Session, request string) (*http.Response, error) {
	payload := `{"id":"` + u.Username + `","userId":` + to_string(u.Id) + `,"username":"` + u.Username + `","thumbnail_url":"//uploads.scratch.mit.edu/users/avatars/` + to_string(u.Id) + `.png"}`

	req, err := http.NewRequest("PUT", "https://scratch.mit.edu/site-api/users/followers/"+u.Username+"/"+request+"/?usernames="+session.Username, bytes.NewBuffer([]byte(payload)))
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

func GetUser(username string) *User {
	resp, err := http.Get("https://api.scratch.mit.edu/users/" + username)
	if err != nil {
		panic(err)
	}

	responseUser := new(User)

	json.NewDecoder(resp.Body).Decode(responseUser)

	return responseUser
}
