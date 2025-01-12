package main

import (
	"scratchgonnect"
)

func main() {
	//newUser := get_user("kajtolmation")
	//newUser.get_following()

	//newStudio := get_studio("31659696")
	//curators := newStudio.get_managers()
	//fmt.Println(curators)

	//newProject := get_project("23728")
	//fmt.Println(newProject.Name

	//go run ./example

	session := scratchgonnect.NewSession("", "")

	project := scratchgonnect.GetProject("535962801")
	project.SetProject(*session, "bot verify1", "Hello", project.Description)

}
