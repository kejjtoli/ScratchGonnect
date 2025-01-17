package main

import (
	"fmt"

	scratchgonnect "github.com/kejjtoli/ScratchGonnect"
)

func start() {
	// Log into scratch account (required for some functions)
	session := scratchgonnect.NewSession("username", "pass")

	// User structure
	user := scratchgonnect.GetUser("kajtolmation")
	user.GetFollowers() // Returns an array of users
	user.GetFollowing() // Returns an array of users
	user.GetComments(1)

	// Authentication required

	user.Follow(*session)                         // Follows user
	user.Unfollow(*session)                       // Unfollows user
	user.PostComment(*session, "content", "", "") // Posts comment under profile

	// Studio structure
	studio := scratchgonnect.GetStudio("34645019")
	studio.GetCurators()      // Returns an array of users
	studio.GetManagers()      // Returns an array of users
	studio.GetComments(0, 20) // Returns list of comments
	studio.GetProjects()      // Returns list of projects

	// Authentication required

	studio.Follow(*session)                     // Follows studio
	studio.Unfollow(*session)                   // Unfollows studio
	studio.AddProject(*session, 535962801)      // Adds project to studio
	studio.RemoveProject(*session, 535962801)   // Removes project from studio
	studio.InviteUser(*session, "kajtolmation") // Invites user to studio

	// Project structure
	project := scratchgonnect.GetProject("535962801")
	project = scratchgonnect.GetHiddenProject(*session, "535962801") // For unshared projects, also works on shared projects
	project.GetRemixes()
	project.GetComments(0, 20)

	// Authentication required

	project.SetProject(*session, "title", "", "")    // Sets project data
	project.PostComment(*session, "content", "", "") // Posts comment under project
	project.Love(*session)
	project.Favorite(*session)
	project.Unlove(*session)
	project.Unfavorite(*session)
	project.Unshare(*session)
	project.Share(*session)

	// Explore page
	scratchgonnect.GetExploreStudios("trending", "en") // Returns list of studios

	// Cloud (Turbowarp)
	cloud := scratchgonnect.ConnectTurbowarpCloud("username", "1121839236") // Connects to turbowarp websocket and performs handshake
	cloud.SetVariable("t1", 314)                                            // Sets cloud variable to value
	cloud.GetVariable("t1")                                                 // Returns value of cloud variable
	cloud.DisconnectCloud()                                                 // Disconnet from websocket
	cloud.Listen(cloud_listener)                                            // Listen for cloud messages
}

func cloud_listener(connection *scratchgonnect.CloudSocket, variable_name string, value int) {
	fmt.Println(variable_name, value)
}
