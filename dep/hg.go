package dep

// Copyright 2013 Vubeology, Inc.

import (
	"github.com/vube/depman/colors"
	"github.com/vube/depman/util"
	"os/exec"
	"strings"
)

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

func (h *Hg) Clone(d *Dependency) (err error) {
	if !util.Exists(d.Path()) {
		err = util.RunCommand("go get -u " + d.Repo)
	}
	return
}

func (h *Hg) Fetch(d *Dependency) (err error) {
	err = util.RunCommand("hg pull")
	return
}

func (h *Hg) Update(d *Dependency) (err error) {
	err = util.RunCommand("hg up" + d.Version)
	return
}

func (h *Hg) Checkout(d *Dependency) (err error) {
	err = util.RunCommand("hg up " + d.Version)
	return
}

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

	{
		var out_bytes []byte
		out_bytes, err = exec.Command("hg", "id", "-i", d.Version).CombinedOutput()
		hash = strings.TrimSuffix(string(out_bytes), "\n")
	}

	util.Cd(pwd)

	if err != nil {
		util.Print("pwd: " + util.Pwd())
		util.PrintIndent(colors.Red("hg id -i " + d.Version))
		util.PrintIndent(colors.Red(hash))
		util.PrintIndent(colors.Red(err.Error()))
		util.Fatal("")
	}

	return
}
