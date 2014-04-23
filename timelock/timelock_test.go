package timelock

// Copyright 2013-2014 Vubeology, Inc.

import (
	"github.com/vube/depman/dep"
	. "launchpad.net/gocheck"
	"testing"
	"time"
)

// Hook up gocheck into the "go test" runner.
func TestTimelock(t *testing.T) {
	TestingT(t)
}

type TimelockSuite struct{}

var _ = Suite(&TimelockSuite{})

func (s *TimelockSuite) TestIsStale(c *C) {
	cache = make(map[string]time.Time)

	old := new(dep.Dependency)
	old.Repo = "old"

	new := new(dep.Dependency)
	new.Repo = "new"

	cache["old"] = time.Now().Add(-2 * time.Hour)
	cache["new"] = time.Now().Add(-1 * time.Minute)

	c.Check(IsStale(old), Equals, true)
	c.Check(IsStale(new), Equals, false)

}
