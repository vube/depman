// Package dep provides a  Dependency struct, a DependencyMap from nicknames (strings) to Dependencies,
// and functions to read and write a DependencyMap to a deps.json file
package dep

// Copyright 2013-2014 Vubeology, Inc.

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/vube/depman/colors"
	"github.com/vube/depman/util"
)

// Dependency Types
const (
	TypeGit      = "git"
	TypeHg       = "hg"
	TypeBzr      = "bzr"
	TypeGitClone = "git-clone"
)

var (
	// ErrUnknownType indicates that an unknown dependency type was found
	ErrUnknownType = errors.New("unknown dependency type")

	// ErrMissingAlias indicates that a git-clone dependency requires an alias field
	ErrMissingAlias = errors.New("dependency type git-clone requires alias field")
)

// DepsFile is the name of the dependency file
const DepsFile string = "deps.json"

// Dependency defines a single dependency
type Dependency struct {
	Repo      string         `json:"repo"`
	Version   string         `json:"version"`
	Type      string         `json:"type"`
	Alias     string         `json:"alias,omitempty"`
	SkipCache bool           `json:"skip-cache,omitempty"`
	VCS       VersionControl `json:"-"`
}

// VersionControl is an interface that define a standard set of operations that can be completed by a version control system
type VersionControl interface {
	Clone(d *Dependency) (err error)

	// Get changes from the server
	Fetch(d *Dependency) (err error)

	// Pull/Merge/Update this branch
	Update(d *Dependency) (err error)

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
	for key := range deps.Map {
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

// SetupVCS configures the VCS depending on the type
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

	if d.Type != TypeGitClone && d.Alias != "" {
		util.Print(colors.Yellow("Warning: " + d.Repo + ": 'alias' field only allowed in dependencies with type 'git-clone', skipping..."))
		d.Alias = ""
	}

	return
}

// Path returns the path for this dependency
// searches for the appropriate directory in each part of the GOPATH (delimited by ':')
// if not found return the path using the first port of GOPATH
func (d *Dependency) Path() (p string) {
	parts := strings.Split(os.Getenv("GOPATH"), ":")

	for _, path := range parts {
		p = filepath.Join(path, "src")
		if d.Alias == "" {
			p = filepath.Join(p, d.Repo)
		} else {
			p = filepath.Join(p, d.Alias)
		}

		if util.Exists(p) {
			return
		}
	}

	// didn't find a directory, use the first part of gopath
	p = filepath.Join(parts[0], "src")
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
