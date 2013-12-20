// Package dep provides a  Dependency struct, a DependencyMap from nicknames (strings) to Dependencies,
// and functions to read and write a DependencyMap to a deps.json file
package dep

// Copyright 2013 Vubeology, Inc.

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"vube/depman/colors"
)

// Dependency Types
const (
	TypeGit      = "git"
	TypeHg       = "hg"
	TypeBzr      = "bzr"
	TypeGitClone = "git-clone"
)

// The name of the dependency file
const DepsFile string = "deps.json"

// Dependency defines a single dependency
type Dependency struct {
	Repo    string `json:"repo"`
	Version string `json:"version,omitempty"`
	Type    string `json:"type"`
	Alias   string `json:"alias,omitempty"`
}

// DependencyMap defines a set of dependencies
type DependencyMap struct {
	Map  map[string]Dependency
	Path string
}

// New returns a newly constructed DependencyMap
func New() (d DependencyMap) {
	d.Map = make(map[string]Dependency)
	return
}

// Read reads filename and parses the content into a DependencyMap
func Read(filename string) (deps DependencyMap, err error) {
	deps.Map = make(map[string]Dependency)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &deps.Map)
	if err != nil {
		return
	}

	for key, _ := range deps.Map {
		val := deps.Map[key]
		// if no version specified, use a default of master
		if val.Version == "" {
			val.Version = "master"
			deps.Map[key] = val
		}
	}

	deps.Path = filename
	return
}

// Write the dependencyMap back into to the file it was read from
func (d *DependencyMap) Write() (err error) {

	var buf bytes.Buffer
	str, err := json.Marshal(d.Map)
	json.Indent(&buf, str, "", "    ")

	if err == nil {
		data := []byte(buf.String() + "\n")
		ioutil.WriteFile(d.Path, data, 0644)
	}
	return
}

// Path returns the path to the deps.json file that this DependencyMap was read from
func (d *Dependency) Path() (p string) {
	goPath := os.Getenv("GOPATH")
	if strings.TrimSpace(goPath) == "" {
		log.Fatal(colors.Red("You must set GOPATH"))
	}

	p = path.Join(goPath, "src")

	if d.Alias == "" {
		p = path.Join(p, d.Repo)
	} else {
		p = path.Join(p, d.Alias)
	}

	return

}
