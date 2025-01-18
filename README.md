<img align="center" src="https://github.com/kejjtoli/ScratchGonnect/blob/main/example/logo.png?raw=true" alt="ScratchGonnect">
 
![Static Badge](https://img.shields.io/badge/License-GPL3.0-blue)
![Commits Since Latest Version](https://img.shields.io/github/commits-since/kejjtoli/ScratchGonnect/latest)

**Send Scratch API requests with ease.**

ScratchGonnect is a Scratch API wrapper written in Go and allows you to easily and quickly interface with the Scratch API. It allows you to post comments, change project data and much more.

## Installation

**Run this command in your terminal in order to use this package**
```
go get github.com/kejjtoli/ScratchGonnect@latest
```

## Examples

**Creating a new session (lets you use features that require logging in)**
```go
// Log into a scratch session
session := scratchgonnect.NewSession("username", "pass")
```
Posting a comment under a project
```go
// Get a project by id
project := scratchgonnect.GetProject("535962801")

// Post comment with given content, parent_id, commentee_id
project.PostComment(*session, "Comment Content", "", "")
```
Following a user and getting his followers
```go
// Get a user by username
user := scratchgonnect.GetUser("kajtolmation")

// Follow user
user.Follow(*session)

// Get array of user objects
followers := user.GetFollowers()
```
**Connecting to Turbowarp cloud**
```go
// Connect to turbowarp websocket
cloud := scratchgonnect.ConnectTurbowarpCloud("username", "1121839236")
```
Setting a cloud variable
```go
// Sets the value of a cloud variable
cloud.SetVariable("t1", 314)
```
Listening to cloud variable changes
```go
// Listens to all set variable messages
cloud.Listen(cloud_listener)

// Prints out name and new value of changed variable
func cloud_listener(variable_name string, value int) {
	fmt.Println(variable_name, value)
}
```

*More examples in examples/example.go file*

## License
ScratchGonnect is licensed under the GPL-3.0 license and is free to use and open-source.
