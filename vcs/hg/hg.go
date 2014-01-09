package hg

// Copyright 2013 Vubeology, Inc.

import (
	"os/exec"
	"strings"

	"github.com/vube/depman/colors"
	"github.com/vube/depman/dep"
	"github.com/vube/depman/util"
)

// LastCommit retrieves the version number of the last commit on branch
// Assumes that the current working directory is in the hg repo
func LastCommit(d dep.Dependency, branch string) (hash string) {
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

//GetHead - Render a revspec to a commit ID
func GetHead(d dep.Dependency) (hash string, err error) {
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
