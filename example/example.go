package main

import scratchgonnect "github.com/kejjtoli/ScratchGonnect"

func start() {
	// Log into scratch account (required for some functions)
	session := scratchgonnect.NewSession("bob_budowlaniec", "bobbudulec")

	// User structure
	user := scratchgonnect.GetUser("kajtolmation")
	user.GetFollowers() // Returns an array of users
	user.GetFollowing() // Returns an array of users

	// Authentication required

	user.Follow(*session)                         // Follows user
	user.Unfollow(*session)                       // Unfollows user
	user.PostComment(*session, "content", "", "") // Posts comment under profile

	// Studio structure
	studio := scratchgonnect.GetStudio("34645019")
	studio.GetCurators() // Returns an array of users
	studio.GetManagers() // Returns an array of users

	// Authentication required

	studio.Follow(*session)   // Follows studio
	studio.Unfollow(*session) // Unfollows studio

	// Project structure
	project := scratchgonnect.GetProject("535962801")

	// Authentication required

	project.SetProject(*session, "title", "", "")    // Sets project data
	project.PostComment(*session, "content", "", "") // Posts comment under project
}
