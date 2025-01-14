package scratchgonnect

import (
	"bytes"
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

func (s Studio) GetCurators() []User {
	resp, err := http.Get("https://api.scratch.mit.edu/studios/" + to_string(s.Id) + "/curators")
	if err != nil {
		panic(err)
	}

	decoded := []User{}

	json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded
}

func (s Studio) GetManagers() []User {
	resp, err := http.Get("https://api.scratch.mit.edu/studios/" + to_string(s.Id) + "/managers")
	if err != nil {
		panic(err)
	}

	decoded := []User{}

	json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded
}

func (s Studio) GetComments(offset int, limit int) []Comment {
	resp, err := http.Get("https://api.scratch.mit.edu/studios/" + to_string(s.Id) + "/comments?offset=" + to_string(offset) + "&limit=" + to_string(limit))
	if err != nil {
		panic(err)
	}

	decoded := []Comment{}

	json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded
}

func (s Studio) GetProjects() []Project {
	resp, err := http.Get("https://api.scratch.mit.edu/studios/" + to_string(s.Id) + "/projects")
	if err != nil {
		panic(err)
	}

	decoded := []project_partial{}

	json.NewDecoder(resp.Body).Decode(&decoded)

	from_partial := []Project{}

	for _, element := range decoded {
		new_element := Project{
			Author: User{
				Username: element.Username,
				Id:       element.CreatorId,
			},
			Id:   element.Id,
			Name: element.Title,
		}

		from_partial = append(from_partial, new_element)
	}

	return from_partial
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

func (s Studio) InviteUser(session Session, username string) {
	req, err := http.NewRequest("PUT", "https://scratch.mit.edu/site-api/users/curators-in/"+to_string(s.Id)+"/invite_curator/?usernames="+username, *new(io.Reader))
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
		panic("Invite user failed! http response:" + to_string(resp.StatusCode))
	}
}

func (s Studio) AddProject(session Session, project_id int) {
	resp, err := change_project_request_studio(s, session, project_id, "POST")

	if err != nil || resp.StatusCode != 200 {
		panic("Add project failed! http response:" + to_string(resp.StatusCode))
	}
}

func (s Studio) RemoveProject(session Session, project_id int) {
	resp, err := change_project_request_studio(s, session, project_id, "DELETE")

	if err != nil || resp.StatusCode != 200 {
		panic("Remove project failed! http response:" + to_string(resp.StatusCode))
	}
}

func (studio Studio) PostComment(session Session, content string, parent_id string, commentee_id string) {
	b, _ := json.Marshal(json_comment{
		Content:  content,
		ParentId: parent_id,
		Id:       commentee_id,
	})

	req, err := http.NewRequest("POST", "https://api.scratch.mit.edu/proxy/comments/studio/"+to_string(studio.Id), bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}

	req.Header = session.HttpHeader

	req.Header.Set("referer", "https://scratch.mit.edu/studios/"+to_string(studio.Id)+"/")
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	req.AddCookie(&session.Cookie)
	req.AddCookie(&csrfCookieDefault)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		panic("Post comment failed! http response:" + to_string(resp.StatusCode))
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

func change_project_request_studio(s Studio, session Session, project_id int, rq string) (*http.Response, error) {
	req, err := http.NewRequest(rq, "https://api.scratch.mit.edu/studios/"+to_string(s.Id)+"/project/"+to_string(project_id), *new(io.Reader))
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
