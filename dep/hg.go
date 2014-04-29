package dep

// Copyright 2013-2014 Vubeology, Inc.

import (
	"os/exec"
	"strings"

	"github.com/vube/depman/colors"
	"github.com/vube/depman/util"
)

// Hg implements the VersionControl interface by using Mercurial
type Hg struct{}

// LastCommit retrieves the version number of the last commit on branch
// Assumes that the current working directory is in the hg repo
func (h *Hg) LastCommit(d *Dependency, branch string) (hash string, err error) {
	c := exec.Command("hg", "log", "--template='{node}\n'", "--limit=1")
	out, err := c.CombinedOutput()

	if err != nil {
		util.Print("pwd: " + util.Pwd())
		util.PrintIndent(colors.Red("hg log --template='{node}\n' --limit=1"))
		util.PrintIndent(colors.Red(string(out)))
		util.PrintIndent(colors.Red(err.Error()))
		util.Fatal("")
	}

	hash = strings.Replace(string(out), "\n", "", -1)
	return
}

// Clone uses go get to clone a mercurial repo
func (h *Hg) Clone(d *Dependency) (err error) {
	if !util.Exists(d.Path()) {
		err = util.RunCommand("go get -u " + d.Repo)
	}
	return
}

// Fetch fetches a mercurial repo
func (h *Hg) Fetch(d *Dependency) (err error) {
	err = util.RunCommand("hg pull")
	return
}

// Update updates a mercurial repo
func (h *Hg) Update(d *Dependency) (err error) {
	err = util.RunCommand("hg up " + d.Version)
	return
}

// Checkout updates a mercurial repo
func (h *Hg) Checkout(d *Dependency) (err error) {
	err = util.RunCommand("hg up " + d.Version)
	return
}

//Clean cleans a mercurial repo
func (h *Hg) Clean(d *Dependency) {
	util.PrintIndent(colors.Red("Cleaning:") + colors.Blue(" hg up --clean "+d.Version))
	util.RunCommand("hg up --clean " + d.Version)
	return
}

//GetHead - Render a revspec to a commit ID
func (h *Hg) GetHead(d *Dependency) (hash string, err error) {
	var pwd string

	pwd = util.Pwd()
	util.Cd(d.Path())
	defer util.Cd(pwd)

	out, err := exec.Command("hg", "id", "-i").CombinedOutput()
	hash = strings.TrimSuffix(string(out), "\n")

	if err != nil {
		util.Print("pwd: " + util.Pwd())
		util.PrintIndent(colors.Red("hg id -i " + d.Version))
		util.PrintIndent(colors.Red(hash))
		util.PrintIndent(colors.Red(err.Error()))
		util.Fatal("")
	}

	return
}
