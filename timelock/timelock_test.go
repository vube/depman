package timelock

// Copyright 2013 Vubeology, Inc.

import (
	"testing"
	"time"
	. "launchpad.net/gocheck"
)

// Hook up gocheck into the "go test" runner.
func TestTimelock(t *testing.T) {
	TestingT(t)
}

type TimelockSuite struct{}

var _ = Suite(&TimelockSuite{})

func (s *TimelockSuite) TestIsStale(c *C) {
	cache = make(map[string]time.Time)

	cache["old"] = time.Now().Add(-2 * time.Hour)
	cache["new"] = time.Now().Add(-1 * time.Minute)

	c.Check(IsStale("old"), Equals, true)
	c.Check(IsStale("new"), Equals, false)

}
