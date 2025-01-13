<img align="center" src="https://github.com/kejjtoli/ScratchGonnect/blob/main/example/logo.png?raw=true" alt="ScratchGonnect">
 
![Static Badge](https://img.shields.io/badge/License-GPL3.0-blue)
![Commits Since Latest Version](https://img.shields.io/github/commits-since/kejjtoli/ScratchGonnect/latest)

**Send Scratch API requests with ease.**

ScratchGonnect is a Scratch API wrapper written in Go and allows you to easily and quickly interface with the Scratch API. It allows you to post comments, change project data and much more.

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

*More examples in examples/example.go file*

## License
ScratchGonnect is licensed under the GPL-3.0 license and is free to use and open-source.