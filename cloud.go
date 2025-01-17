package scratchgonnect

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	websocket "github.com/gorilla/websocket"
)

// Structs

type handshake_cloud struct {
	Method    string `json:"method"`
	Username  string `json:"user"`
	ProjectId string `json:"project_id"`
}

type set_cloud struct {
	Method    string `json:"method"`
	Username  string `json:"user"`
	ProjectId string `json:"project_id"`
	Name      string `json:"name"`
	Value     int    `json:"value"`
}

type CloudSocket struct {
	Connection *websocket.Conn
	ProjectId  string
	Variables  []CloudVariable
}

type CloudVariable struct {
	Name  string
	Value int
}

var default_timeout time.Duration = 5000000000

// Struct functions

func ConnectHandshake(s Session, project_id string) *CloudSocket {
	header := http.Header{}
	//header.Add("Cookie", "scratchsessionsid="+s.SessionId+";")
	header.Set("Origin", "https://turbowarp.org")
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36")
	header.Set("Host", "clouddata.turbowarp.org")

	c, _, err := websocket.DefaultDialer.Dial("wss://clouddata.turbowarp.org/", header)
	if err != nil {
		panic(err)
	}

	new_socket := CloudSocket{
		Connection: c,
		ProjectId:  project_id,
	}

	c.SetWriteDeadline(time.Now().Add(default_timeout))

	err = c.WriteJSON(handshake_cloud{
		Method:    "handshake",
		Username:  s.Username,
		ProjectId: project_id,
	})
	if err != nil {
		panic(err)
	}

	c.SetReadDeadline(time.Now().Add(default_timeout))

	_, p, err := c.ReadMessage()
	if err != nil {
		panic(err)
	}

	sets := strings.Split(string(p), "\n")

	for _, set_action := range sets {
		action_decoded := set_cloud{}
		err := json.Unmarshal([]byte(set_action), &action_decoded)
		if err != nil {
			panic(err)
		}

		if action_decoded.Method == "set" {
			new_socket.Variables = append(new_socket.Variables, CloudVariable{
				Name:  action_decoded.Name,
				Value: action_decoded.Value,
			})
		}
	}

	//c.Close()

	return &new_socket
}

func (c CloudSocket) SetVariable(s Session, varname string, value int) {
	err := c.Connection.WriteJSON(set_cloud{
		Method:    "set",
		Username:  s.Username,
		ProjectId: c.ProjectId,
		Name:      "☁ " + varname,
		Value:     value,
	})
	if err != nil {
		panic(err)
	}

	for _, variable := range c.Variables {
		if variable.Name == "☁ "+varname {
			variable.Value = value
		}
	}
}
