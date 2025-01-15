package scratchgonnect

import (
	"encoding/json"
	"net/http"
)

// Funcs

func GetExploreStudios(mode string, language string) []Studio {
	resp, err := http.Get("https://api.scratch.mit.edu/explore/studios?q=animations&mode=" + mode + "&language=" + language)
	if err != nil {
		panic(err)
	}

	decoded := []Studio{}

	json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded
}
