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
	Username   string
}

type CloudVariable struct {
	Name  string
	Value int
}

type fn func(*CloudSocket, string, int)

var default_timeout time.Duration = 5000000000

// Struct functions

func (c CloudSocket) DisconnectCloud() {
	// Closes websocket connection
	c.Connection.Close()
}

func (c CloudSocket) SetVariable(varname string, value int) {
	// Sets a cloud variable
	if get_variable_key(c.Variables, varname) != value {
		c.Connection.SetWriteDeadline(time.Now().Add(default_timeout))

		err := c.Connection.WriteJSON(set_cloud{
			Method:    "set",
			Username:  c.Username,
			ProjectId: c.ProjectId,
			Name:      "☁ " + varname,
			Value:     value,
		})
		if err != nil {
			panic(err)
		}

		for i, variable := range c.Variables {
			if variable.Name == "☁ "+varname {
				c.Variables[i].Value = value
				break
			}
		}
	}
}

func (c CloudSocket) GetVariable(varname string) int {
	// Get variable value from CloudSocket
	for _, variable := range c.Variables {
		if variable.Name == "☁ "+varname {
			return variable.Value
		}
	}

	panic("Couldn't find variable!")
}

func (c CloudSocket) Listen(f fn) {
	// Perform infinite loop for handling cloud events
	for {
		c.Connection.SetReadDeadline(time.Time{})

		// Wait for a message to get recieved
		_, p, err := c.Connection.ReadMessage()
		if err != nil {
			panic(err)
		}

		// Update cloud variable values and fire listener function
		sets := strings.Split(string(p), "\n")

		for _, set_action := range sets {
			action_decoded := set_cloud{}
			err := json.Unmarshal([]byte(set_action), &action_decoded)
			if err != nil {
				panic(err)
			}

			if action_decoded.Method == "set" {
				for i, cloud_var := range c.Variables {
					if cloud_var.Name == action_decoded.Name {
						c.Variables[i].Value = action_decoded.Value
						f(&c, cloud_var.Name, action_decoded.Value)
					}
				}
			}
		}
	}
}

// Functions

func get_variable_key(variables []CloudVariable, name string) int {
	// Search for variable name in CloudVariable array and return its value
	for _, variable := range variables {
		if variable.Name == "☁ "+name {
			return variable.Value
		}
	}

	return 0
}

func ConnectTurbowarpCloud(username string, project_id string) *CloudSocket {
	header := http.Header{}
	//header.Add("Cookie", "scratchsessionsid="+s.SessionId+";")
	header.Set("Origin", "https://turbowarp.org")
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36")
	header.Set("Host", "clouddata.turbowarp.org")

	// Connect to Turbowarp cloud websocket
	c, _, err := websocket.DefaultDialer.Dial("wss://clouddata.turbowarp.org/", header)
	if err != nil {
		panic(err)
	}

	// If succesfull make a new CloudSocket object
	new_socket := CloudSocket{
		Connection: c,
		ProjectId:  project_id,
		Username:   username,
	}

	c.SetWriteDeadline(time.Now().Add(default_timeout))

	// Perform handshake with Turbowarp cloud websocket
	err = c.WriteJSON(handshake_cloud{
		Method:    "handshake",
		Username:  username,
		ProjectId: project_id,
	})
	if err != nil {
		panic(err)
	}

	c.SetReadDeadline(time.Now().Add(default_timeout))

	// Receive information about current cloud data values
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

	return &new_socket
}
