// Dependency Manager for Golang Projects
// Author: Nicholas Capo <nicholas@vubeology.com>
//
// For help run 'depman help'
package main

// Copyright 2013 Vubeology, Inc.

//===============================================

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"vube/depman/add"
	"vube/depman/colors"
	"vube/depman/create"
	"vube/depman/dep"
	"vube/depman/install"
	"vube/depman/showfrozen"
	"vube/depman/update"
	"vube/depman/upgrade"
	"vube/depman/util"
)

// Version number
const VERSION string = "2.5.2"

//===============================================

func main() {
	var result int
	var help bool
	var path string
	var command string
	var arguments []string
	var deps dep.DependencyMap
	var err error

	log.SetFlags(0)

	flag.BoolVar(&help, "help", false, "Display help")
	flag.StringVar(&path, "path", ".", "Directory or full path to deps.json")
	util.Parse()

	util.Version(VERSION)

	path = util.GetPath(path)

	if flag.NArg() > 0 {
		command = strings.ToLower(flag.Arg(0))
	}

	if flag.NArg() > 1 {
		arguments = flag.Args()[1:]
	}

	if help {
		command = "help"
	}

	// switch to check for deps.json
	switch command {
	case "init", "help":
		// don't check for deps.json
	default:
		util.CheckPath(path)
		deps, err = dep.Read(path)
		if err != nil {
			util.Fatal(colors.Red("Error Reading deps.json: " + err.Error()))
		}
	}

	// switch to exec the sub command
	switch command {
	case "init", "create":
		result = create.Create(path)
	case "add":
		if len(arguments) < 1 {
			util.Print(colors.Red("Add command requires 1 argument: Add [nickname]"))
			Help()
		} else {
			add.Add(deps, arguments[0])
		}

	case "update":
		if len(arguments) < 2 {
			util.Print(colors.Red("Update command requires 2 arguments: Update [nickname] [branch]"))
			Help()
		} else {
			update.Update(deps, arguments[0], arguments[1])
		}
	case "install", "":
		result = install.Install(deps)
	case "self-upgrade":
		upgrade.Self()
	case "show-frozen":
		var recursive bool
		flagset := flag.NewFlagSet("show-frozen", flag.ExitOnError)
		flagset.BoolVar(&recursive, "recursive", false, "descend recursively (depth-first) into dependencies")
		flagset.Parse(flag.Args()[1:])

		if recursive {
			fmt.Println(showfrozen.ReadRecursively(deps, nil))
		} else {
			fmt.Print(showfrozen.Read(deps))
		}
	default:
		result = 1
		log.Println(colors.Red("Unknown Command: " + command))
		fallthrough
	case "help":
		Help()
	}

	if result == 0 {
		util.Print("Success")
	}

	util.OsExit(result)
}

//===============================================

// Help prints the help message for depman
func Help() {
	log.Println("")
	log.Println("Commands:")
	log.Println("   Init                        : Create an empty deps.json (not yet implemented)")
	log.Println("   Add [nickname]              : Add a dependency (interactive)")
	log.Println("   Install                     : Install all the dependencies listed in deps.json (default)")
	log.Println("   Update [nickname] [branch]  : Update [nickname] to use the latest commit in [branch]")
	log.Println("   Self-Upgrade                : Upgrade depman to the latest version on the master branch")
	log.Println("   Help                        : Display this help")
	log.Println("   show-frozen                 : Show dependencies as resolved to commit IDs")
	log.Println("")
	log.Println("Example: depman --verbose install")
	log.Println("")
	//log.Println("   freeze                      : For each dependency change tag and branch versions to commits (not yet implemented)")
	log.Println("Options:")
	flag.PrintDefaults()
}
