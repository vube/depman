package install

// Copyright 2013 Vubeology, Inc.

import (
	"bytes"
	"log"
	"testing"

	. "launchpad.net/gocheck"

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
	clean = false
}

/*
func (s *TestSuite) TestGitClone(c *C) {
	var cmds []string
	var b []byte

	dir, err := ioutil.TempDir("", "DepmanUnitTest")
	c.Check(err, IsNil)
	os.Setenv("GOPATH", dir)
	os.Mkdir(dir+"/src", 0755)
	defer os.RemoveAll(dir)

	util.RunCommand = func(s string) int {
		cmds = append(cmds, s)
		return 0
	}

	buf := bytes.NewBuffer(b)
	log.SetOutput(buf)

	d := dep.Dependency{Repo: "REPO"}

	GitClone(d, dir+"/src/REPO")
	c.Check(cmds[0], Equals, "git clone REPO "+dir+"/src/REPO")
	c.Check(s.buf.String(), Equals, "")
	cmds = []string{}

	os.MkdirAll(dir+"/src/REPO", 0755)
	GitClone(d, dir+"/src/REPO")
	c.Check(cmds[0], Equals, "git fetch")
	c.Check(util.Pwd(), Equals, dir+"/src/REPO")
}
*/

func (s *TestSuite) TestDuplicate(c *C) {
	set := make(map[string]string)

	d := dep.Dependency{Repo: "repo", Version: "version", Type: "type"}
	c.Check(len(set), Equals, 0)

	// No dup
	skip := duplicate(d, set)
	c.Check(skip, Equals, false)
	c.Check(len(set), Equals, 1)
	v, ok := set["repo"]
	c.Check(ok, Equals, true)
	c.Check(v, Equals, "version")

	// dup same version
	util.SetVerbose(true)
	skip = duplicate(d, set)
	c.Check(len(set), Equals, 1)
	c.Check(s.buf.String(), Equals, "Skipping previously installed dependency: repo\n")
	c.Check(skip, Equals, true)

	s.buf.Truncate(0)

	// dup different version
	d.Version = "version2"
	skip = duplicate(d, set)
	c.Check(len(set), Equals, 1)
	out := "ERROR    : Duplicate dependency with different versions detected\nRepo     : repo\nVersions : version2\tversion\n"
	c.Check(s.buf.String(), Equals, out)
	c.Check(skip, Equals, false)

}
