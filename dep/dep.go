// Package dep provides a  Dependency struct, a DependencyMap from nicknames (strings) to Dependencies,
// and functions to read and write a DependencyMap to a deps.json file
package dep

// Copyright 2013 Vubeology, Inc.

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/vube/depman/colors"
	"github.com/vube/depman/util"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Dependency Types
const (
	TypeGit      = "git"
	TypeHg       = "hg"
	TypeBzr      = "bzr"
	TypeGitClone = "git-clone"
)

var (
	ErrUnknownType  = errors.New("unknown dependency type")
	ErrMissingAlias = errors.New("dependency type git-clone requires alias field")
)

// The name of the dependency file
const DepsFile string = "deps.json"

// Dependency defines a single dependency
type Dependency struct {
	Repo    string         `json:"repo"`
	Version string         `json:"version"`
	Type    string         `json:"type"`
	Alias   string         `json:"alias,omitempty"`
	VCS     VersionControl `json:"-"`
}

type VersionControl interface {
	Clone(d *Dependency) (err error)
	Pull(d *Dependency) (err error)
	Checkout(d *Dependency) (err error)

	LastCommit(d *Dependency, branch string) (hash string, err error)
	GetHead(d *Dependency) (to_return string, err error)

	Clean(d *Dependency)
}

// DependencyMap defines a set of dependencies
type DependencyMap struct {
	Map  map[string]*Dependency
	Path string
}

// New returns a newly constructed DependencyMap
func New() (d DependencyMap) {
	d.Map = make(map[string]*Dependency)
	return
}

// Read reads filename and parses the content into a DependencyMap
func Read(filename string) (deps DependencyMap, err error) {
	deps.Map = make(map[string]*Dependency)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &deps.Map)
	if err != nil {
		return
	}

	// traverse map and look for empty version fields - provide a default if such found
	for key, _ := range deps.Map {
		val := deps.Map[key]
		if val.Version == "" {
			switch val.Type {
			case TypeGit, TypeGitClone:
				val.Version = "master"
			case TypeHg:
				val.Version = "tip"
			case TypeBzr:
				val.Version = "trunk"
			default:
				val.Version = ""
			}
			deps.Map[key] = val
		}
	}

	for name, d := range deps.Map {
		err := d.SetupVCS(name)
		if err != nil {
			delete(deps.Map, name)
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

// Configures the VCS depending on the type
func (d *Dependency) SetupVCS(name string) (err error) {
	switch d.Type {
	case TypeGitClone:
		if d.Alias == "" {
			util.PrintIndent(colors.Red("Error: Dependency " + name + ": Repo '" + d.Repo + "' Type '" + d.Type + "' requires 'alias' field"))
			err = ErrMissingAlias
			return
		}

		d.VCS = new(Git)
	case TypeGit:
		d.VCS = new(Git)
	case TypeBzr:
		d.VCS = new(Bzr)
	case TypeHg:
		d.VCS = new(Hg)
	default:
		util.PrintIndent(colors.Red(d.Repo + ": Unknown repository type (" + d.Type + "), skipping..."))
		util.PrintIndent(colors.Red("Valid Repository types: " + TypeGit + ", " + TypeHg + ", " + TypeBzr + ", " + TypeGitClone))
		err = ErrUnknownType
	}

	return
}

// Path returns the path to the deps.json file that this DependencyMap was read from
func (d *Dependency) Path() (p string) {
	goPath := os.Getenv("GOPATH")
	if strings.TrimSpace(goPath) == "" {
		log.Fatal(colors.Red("You must set GOPATH"))
	}

	p = filepath.Join(goPath, "src")

	if d.Alias == "" {
		p = filepath.Join(p, d.Repo)
	} else {
		p = filepath.Join(p, d.Alias)
	}

	return

}

//GetPath processes p and returns a clean path ending in deps.json
func GetPath(p string) (result string) {
	if !strings.HasSuffix(p, DepsFile) {
		result = p + "/" + DepsFile
	}
	result = filepath.Clean(result)
	return
}
