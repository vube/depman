package vcs

// Copyright 2013 Vubeology, Inc.

import (
	"bytes"
	. "launchpad.net/gocheck"
	"log"
	"testing"
	"github.com/vube/depman/colors"
	"github.com/vube/depman/dep"
	"github.com/vube/depman/util"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

type TestSuite struct {
	buf *bytes.Buffer
}

var _ = Suite(&TestSuite{})

func (s *TestSuite) SetUpTest(c *C) {
	colors.Mock()
	var b []byte
	s.buf = bytes.NewBuffer(b)
	util.Mock(s.buf)
	log.SetFlags(0)
}

func (s *TestSuite) TestCheckout(c *C) {

	var cmds []string
	util.RunCommand = func(s string) int {
		cmds = append(cmds, s)
		return 0
	}

	d := dep.Dependency{Repo: "repo", Version: "version"}

	// unkown
	Checkout(d, false)
	out := "repo: Unknown repository type (), skipping...\nValid Repository types: git, hg, bzr, git-clone\n"
	c.Check(s.buf.String(), Equals, out)
	c.Check(len(cmds), Equals, 0)
	s.buf.Truncate(0)
	cmds = []string{}

	// git
	d.Type = "git"
	Checkout(d, false)
	c.Check(len(cmds), Equals, 1)
	c.Check(cmds[0], Equals, "git checkout version")
	s.buf.Truncate(0)
	cmds = []string{}

	// git with clean
	d.Type = "git"
	Checkout(d, true)
	c.Check(len(cmds), Equals, 3)
	c.Check(cmds[0], Equals, "git reset --hard HEAD")
	c.Check(cmds[1], Equals, "git clean -f")
	c.Check(cmds[2], Equals, "git checkout version")
	s.buf.Truncate(0)
	cmds = []string{}

	// hg
	d.Type = "hg"
	Checkout(d, false)
	c.Check(s.buf.String(), Equals, "")
	c.Check(len(cmds), Equals, 1)
	c.Check(cmds[0], Equals, "hg up version")
	s.buf.Truncate(0)
	cmds = []string{}

	// hg with clean
	d.Type = "hg"
	Checkout(d, true)
	out = "Cleaning: hg up --clean version\n"
	c.Check(s.buf.String(), Equals, out)
	c.Check(len(cmds), Equals, 1)
	c.Check(cmds[0], Equals, "hg up --clean version")
	s.buf.Truncate(0)
	cmds = []string{}

	// bzr
	d.Type = "bzr"
	Checkout(d, false)
	c.Check(s.buf.String(), Equals, "")
	c.Check(len(cmds), Equals, 1)
	c.Check(cmds[0], Equals, "bzr up --revision version")
	s.buf.Truncate(0)
	cmds = []string{}

	// bzr with clean
	d.Type = "bzr"
	Checkout(d, true)
	c.Check(s.buf.String(), Equals, "")
	c.Check(len(cmds), Equals, 1)
	c.Check(cmds[0], Equals, "bzr up --revision version")
	s.buf.Truncate(0)
	cmds = []string{}
}
