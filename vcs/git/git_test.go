package git

// Copyright 2013 Vubeology, Inc.

import (
	"bytes"
	. "launchpad.net/gocheck"
	"log"
	"os"
	"testing"
	"github.com/vube/depman/colors"
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
	b := []byte{}
	s.buf = bytes.NewBuffer(b)
	util.Mock(s.buf)
	log.SetFlags(0)
}

func (s *TestSuite) TestIsGitBranch(c *C) {

	util.Cd(os.Getenv("GOPATH") + "/src/github.com/vube/depman")
	c.Check(IsBranch("master"), Equals, true)

	c.Check(IsBranch("2.1.0"), Equals, false)
	c.Check(IsBranch("7da42054c10f55d5f479b84f59013818ccbd1fd7"), Equals, false)

	util.Cd("/")

	c.Check(IsBranch("master"), Equals, false)
	output := "pwd: /\n" +
		"git branch -r\n" +
		"fatal: Not a git repository (or any of the parent directories): .git\n" +
		"exit status 128\n"
	c.Check(s.buf.String(), Equals, output)
}
