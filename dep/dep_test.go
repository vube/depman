package dep

// Copyright 2013 Vubeology, Inc.

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func TestDep(t *testing.T) {
	TestingT(t)
}

type DepSuite struct{}

var _ = Suite(&DepSuite{})

func (s *DepSuite) TestRead(c *C) {

	deps, err := Read("../tests/unit/unit.json")

	c.Assert(err, IsNil)
	c.Assert(len(deps.Map), Equals, 3)

	d, ok := deps.Map["one"]
	c.Check(ok, Equals, true)
	c.Check(d.Repo, Equals, "repo_one")
	c.Check(d.Version, Equals, "1")
	c.Check(d.Type, Equals, "git")

	d, ok = deps.Map["two"]
	c.Check(ok, Equals, true)
	c.Check(d.Repo, Equals, "repo_two")
	c.Check(d.Version, Equals, "2")
	c.Check(d.Type, Equals, "git")

	d, ok = deps.Map["three"]
	c.Check(ok, Equals, true)
	c.Check(d.Repo, Equals, "repo_three")
	c.Check(d.Version, Equals, "3")
	c.Check(d.Type, Equals, "git")

	deps, err = Read("./tests/unit/none")
	c.Check(err, ErrorMatches, "open ./tests/unit/none: no such file or directory")
}
