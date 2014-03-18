package dep

// Copyright 2013 Vubeology, Inc.

import (
	"bytes"
	"log"
	"os"

	"github.com/vube/depman/colors"
	"github.com/vube/depman/util"
	. "launchpad.net/gocheck"
)

type GitSuite struct {
	buf *bytes.Buffer
}

var _ = Suite(&GitSuite{})

func (s *GitSuite) SetUpTest(c *C) {
	colors.Mock()
	b := []byte{}
	s.buf = bytes.NewBuffer(b)
	util.Mock(s.buf)
	log.SetFlags(0)
}

func (s *GitSuite) TestIsGitBranch(c *C) {
	g := new(Git)

	util.Cd(os.Getenv("GOPATH") + "/src/github.com/vube/depman")
	c.Check(g.isBranch("master"), Equals, true)

	c.Check(g.isBranch("2.1.0"), Equals, false)
	c.Check(g.isBranch("7da42054c10f55d5f479b84f59013818ccbd1fd7"), Equals, false)

	util.Cd("/")

	c.Check(g.isBranch("master"), Equals, false)
	output := "pwd: /\n" +
		"git branch -r\n" +
		"fatal: Not a git repository (or any of the parent directories): .git\n" +
		"exit status 128\n"
	c.Check(s.buf.String(), Equals, output)
}
