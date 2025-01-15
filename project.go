package scratchgonnect

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Structs

type Project struct {
	Id              int    `json:"id"`
	Name            string `json:"title"`
	Description     string `json:"description"`
	Instructions    string `json:"instructions"`
	Visiblity       string `json:"visibility"`
	IsPublic        bool   `json:"public"`
	CommentsAllowed bool   `json:"comments_allowed"`
	IsPublished     bool   `json:"is_published"`
	Author          User   `json:"author"`
	Image           string `json:"image"`
	Images          struct {
		Res282x218 string `json:"282x218"`
		Res216x163 string `json:"216x163"`
		Res200x200 string `json:"200x200"`
		Res144x108 string `json:"144x108"`
		Res135x102 string `json:"135x102"`
		Res100x80  string `json:"100x80"`
	} `json:"images"`
	History struct {
		CreatedDate  string `json:"created"`
		ModifiedDate string `json:"modified"`
		SharedDate   string `json:"shared"`
	} `json:"history"`
	Stats struct {
		Views     int `json:"views"`
		Loves     int `json:"loves"`
		Favorites int `json:"favorites"`
		Remixes   int `json:"remixes"`
	} `json:"stats"`
	Remix struct {
		Parent int `json:"parent"`
		Root   int `json:"root"`
	} `json:"remix"`
	Token string `json:"project_token"`
}

type put_http struct {
	Title        string `json:"title"`
	Instructions string `json:"instructions"`
	Description  string `json:"description"`
}

type json_comment struct {
	Content  string `json:"content"`
	ParentId string `json:"parent_id"`
	Id       string `json:"commentee_id"`
}

// Struct functions

func (project Project) SetProject(session Session, title string, instructions string, description string) {
	b, _ := json.Marshal(put_http{
		Title:        title,
		Instructions: instructions,
		Description:  description,
	})

	req, err := http.NewRequest("PUT", "https://api.scratch.mit.edu/projects/"+to_string(project.Id), bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}

	req.Header = session.HttpHeader

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		panic("Set project failed! http response:" + to_string(resp.StatusCode))
	}
}

func (project Project) PostComment(session Session, content string, parent_id string, commentee_id string) {
	b, _ := json.Marshal(json_comment{
		Content:  content,
		ParentId: parent_id,
		Id:       commentee_id,
	})

	req, err := http.NewRequest("POST", "https://api.scratch.mit.edu/proxy/comments/project/"+to_string(project.Id), bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}

	req.Header = session.HttpHeader

	req.Header.Set("referer", "https://scratch.mit.edu/projects/"+to_string(project.Id)+"/")
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	req.AddCookie(&session.Cookie)
	req.AddCookie(&csrfCookieDefault)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		panic("Post comment failed! http response:" + to_string(resp.StatusCode))
	}
}

func (p Project) Love(session Session) {
	resp, err := request_project_post(p, session, "loves", "POST")

	if err != nil || resp.StatusCode != 200 {
		panic("Love project failed! http response:" + to_string(resp.StatusCode))
	}
}

func (p Project) Favorite(session Session) {
	resp, err := request_project_post(p, session, "favorites", "POST")

	if err != nil || resp.StatusCode != 200 {
		panic("Fav project failed! http response:" + to_string(resp.StatusCode))
	}
}

func (p Project) Unlove(session Session) {
	resp, err := request_project_post(p, session, "loves", "DELETE")

	if err != nil || resp.StatusCode != 200 {
		panic("Unlove project failed! http response:" + to_string(resp.StatusCode))
	}
}

func (p Project) Unfavorite(session Session) {
	resp, err := request_project_post(p, session, "favorites", "DELETE")

	if err != nil || resp.StatusCode != 200 {
		panic("Unfav project failed! http response:" + to_string(resp.StatusCode))
	}
}

func (p Project) GetRemixes() []Project {
	resp, err := http.Get("https://api.scratch.mit.edu/projects/" + to_string(p.Id) + "/remixes")
	if err != nil || resp.StatusCode != 200 {
		panic("Get remixes action failed! http response:" + to_string(resp.StatusCode))
	}

	decoded := []Project{}

	json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded
}

// Functions

func request_project_post(p Project, session Session, request string, rq string) (*http.Response, error) {
	req, err := http.NewRequest(rq, "https://api.scratch.mit.edu/proxy/projects/"+to_string(p.Id)+"/"+request+"/user/"+session.Username, *new(io.Reader))
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

func GetProject(project string) *Project {
	resp, err := http.Get("https://api.scratch.mit.edu/projects/" + project)
	if err != nil {
		panic(err)
	}

	response := new(Project)

	json.NewDecoder(resp.Body).Decode(response)

	return response
}
