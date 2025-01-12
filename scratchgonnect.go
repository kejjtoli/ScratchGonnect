package scratchgonnect

import (
	"fmt"
)

// Base structs

type JSONPostLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenDelivery [1]struct {
	Token string `json:"token"`
}

type UserArray []struct {
	User
}

// Functions

func to_string(v int) string {
	return fmt.Sprintf("%d", v)
}
