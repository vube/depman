package dep

// Copyright 2013 Vubeology, Inc.

import (
	"github.com/vube/depman/colors"
	"github.com/vube/depman/util"
	"os/exec"
	"strings"
)

type Bzr struct{}

// LastCommit retrieves the version number of the last commit on branch
// Assumes that the current working directory is in the bzr repo
func (b *Bzr) LastCommit(d *Dependency, branch string) (hash string, err error) {
	c := exec.Command("bzr", "log", "--line")
	out, err := c.CombinedOutput()

	if err != nil {
		util.Print("pwd: " + util.Pwd())
		util.PrintIndent(colors.Red("bzr log --line"))
		util.PrintIndent(colors.Red(string(out)))
		util.PrintIndent(colors.Red(err.Error()))
		util.Fatal("")
	}

	hash = strings.Split(string(out), ":")[0]
	return
}

//GetHead - Render a revspec to a commit ID
func (b *Bzr) GetHead(d *Dependency) (hash string, err error) {
	var pwd string

	pwd = util.Pwd()
	util.Cd(d.Path())

	{
		var out_bytes []byte
		out_bytes, err = exec.Command("bzr", "revno", d.Version).CombinedOutput()
		hash = strings.TrimSuffix(string(out_bytes), "\n")
	}

	util.Cd(pwd)

	if err != nil {
		util.Print("pwd: " + util.Pwd())
		util.PrintIndent(colors.Red("bzr revno " + d.Version))
		util.PrintIndent(colors.Red(hash))
		util.PrintIndent(colors.Red(err.Error()))
		util.Fatal("")
	}

	return
}

func (b *Bzr) Clone(d *Dependency) (err error) {
	if !util.Exists(d.Path()) {
		err = util.RunCommand("go get -u " + d.Repo)
	}
	return
}

func (b *Bzr) Update(d *Dependency) (err error) {
	return
}

func (b *Bzr) Fetch(d *Dependency) (err error) {
	err = util.RunCommand("bzr pull")
	return
}

func (b *Bzr) Checkout(d *Dependency) (err error) {
	err = util.RunCommand("bzr up --revision " + d.Version)
	return
}

func (b *Bzr) Clean(d *Dependency) {
	return
}
