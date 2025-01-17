package scratchgonnect

import (
	"fmt"
	"strings"
)

// Base structs

type jSONPostLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenDelivery [1]struct {
	Token string `json:"token"`
}

type project_partial struct {
	Username  string `json:"username"`
	Id        int    `json:"id"`
	CreatorId int    `json:"creator_id"`
	Title     string `json:"title"`
	History   struct {
		CreatedDate  string `json:"created"`
		ModifiedDate string `json:"modified"`
	} `json:"history"`
}

type json_msgs struct {
	Count int `json:"count"`
}

// Functions

func to_string(v int) string {
	return fmt.Sprintf("%d", v)
}

func scrape_element(el string, target string) string {
	split := strings.Split(el, ` `)
	for _, field := range split {
		args := strings.Split(field, `=`)
		if len(args) == 2 {
			if args[0] == target {
				return strings.Split(args[1], `"`)[1]
			}
		}
	}

	return ""
}
