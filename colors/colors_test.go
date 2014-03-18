package colors

// Copyright 2013 Vubeology, Inc.

import (
	"testing"
	. "launchpad.net/gocheck"
)

// Hook up gocheck into the "go test" runner.
func ColorsTest(t *testing.T) {
	TestingT(t)
}

type ColorsSuite struct{}

var _ = Suite(&ColorsSuite{})

//===============================================

func (s *ColorsSuite) SetUpTest(c *C) {
	noColors = true
}

//===============================================

func (s *ColorsSuite) TestRed(c *C) {
	noColors = false
	str := "NNN"
	r := Red(str)
	c.Check(r, Equals, "\x1b[31m"+str+"\x1b[0m")

	noColors = true

	r = Red(str)
	c.Check(r, Equals, str)
}

func (s *ColorsSuite) TestBlue(c *C) {
	noColors = false
	str := "NNN"
	b := Blue(str)
	c.Check(b, Equals, "\x1b[36m"+str+"\x1b[0m")

	noColors = true

	b = Blue(str)
	c.Check(b, Equals, str)
}

func (s *ColorsSuite) TestYellow(c *C) {
	noColors = false
	str := "NNN"
	y := Yellow(str)
	c.Check(y, Equals, "\x1b[33m"+str+"\x1b[0m")

	noColors = true

	y = Yellow(str)
	c.Check(y, Equals, str)
}
