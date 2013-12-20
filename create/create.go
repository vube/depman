// Package create provides functions to initialize a new deps.json
package create

// Copyright 2013 Vubeology, Inc.

import (
	"io/ioutil"
	"vube/depman/colors"
	"vube/depman/dep"
	"vube/depman/util"
)

// Default empty deps.json
const template = "{}\n"

// Create writes an empty deps.json at the location specified by path
func Create(path string) (result int) {
	if util.Exists(path) {
		result = 1
		util.Fatal(colors.Red(dep.DepsFile + " already exists!"))
	}
	util.Print(colors.Blue("Initializing:"))
	err := ioutil.WriteFile(path, []byte(template), 0644)
	if err == nil {
		util.Print("Empty " + dep.DepsFile + " created (" + path + ")")
	} else {
		util.Fatal(colors.Red("Error creating "+dep.DepsFile+": "), err)
	}
	return
}
