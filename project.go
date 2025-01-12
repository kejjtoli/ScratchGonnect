package scratchgonnect

import (
	"bytes"
	"encoding/json"
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

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
}

// Functions

func GetProject(project string) *Project {
	resp, err := http.Get("https://api.scratch.mit.edu/projects/" + project)
	if err != nil {
		panic(err)
	}

	response := new(Project)

	json.NewDecoder(resp.Body).Decode(response)

	return response
}
