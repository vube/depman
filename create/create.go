// Package create provides functions to initialize a new deps.json
package create

// Copyright 2013-2014 Vubeology, Inc.

import (
	"io/ioutil"

	"github.com/vube/depman/colors"
	"github.com/vube/depman/dep"
	"github.com/vube/depman/util"
)

// Default empty deps.json
const template = "{}\n"

// Create writes an empty deps.json at the location specified by path
func Create(path string) {
	if util.Exists(path) {
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
