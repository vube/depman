package upgrade

// Copyright 2013-2014 Vubeology, Inc.

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func TestUpgrade(t *testing.T) {
	TestingT(t)
}

type UpgradeSuite struct {
	getter func(url string) (data []byte, err error)
}

var _ = Suite(&UpgradeSuite{})

func (s *UpgradeSuite) SetUpSuite(c *C) {
	s.getter = getter
	getter = fromStatic
}

func (s *UpgradeSuite) TearDownSuite(c *C) {
	getter = s.getter
}

func (s *UpgradeSuite) TestVersionParseEmpty(c *C) {
	v := new(version)
	err := v.parse("")
	c.Check(err, Not(Equals), nil)
	c.Check(err.Error(), Equals, "versions must be three ints separated by '.'")
}

func (s *UpgradeSuite) TestVersionParseSimple(c *C) {
	v := new(version)
	err := v.parse("1.2.3")
	c.Assert(err, Equals, nil)
	c.Check(v.a, Equals, 1)
	c.Check(v.b, Equals, 2)
	c.Check(v.c, Equals, 3)
}

func (s *UpgradeSuite) TestVersionParseTwoDigit(c *C) {
	v := new(version)
	err := v.parse("1.10.99")
	c.Assert(err, Equals, nil)
	c.Check(v.a, Equals, 1)
	c.Check(v.b, Equals, 10)
	c.Check(v.c, Equals, 99)

}

func (s *UpgradeSuite) TestVersionParseLeadingZero(c *C) {
	v := new(version)
	err := v.parse("1.02.003")
	c.Assert(err, Equals, nil)
	c.Check(v.a, Equals, 1)
	c.Check(v.b, Equals, 2)
	c.Check(v.c, Equals, 3)

}

func (s *UpgradeSuite) TestCheck(c *C) {
	r, err := check("1.0.0")
	c.Assert(err, Equals, nil)
	c.Check(r, Equals, "2.0.0")

	r, err = check("1.5.0")
	c.Assert(err, Equals, nil)
	c.Check(r, Equals, "2.0.0")

	r, err = check("2.5.0")
	c.Assert(err, Equals, nil)
	c.Check(r, Equals, "")

}

func (s *UpgradeSuite) TestMax(c *C) {

	versions := []*version{
		&version{"1.1.1", 1, 1, 1},
		&version{"2.1.9", 2, 1, 9},
		&version{"8.0.0", 8, 0, 0},
	}
	m := max(versions)
	c.Check(m.str, Equals, "8.0.0")
	c.Check(m.a, Equals, 8)
	c.Check(m.b, Equals, 0)
	c.Check(m.c, Equals, 0)
}

func (s *UpgradeSuite) TestVersionGreaterThan(c *C) {
	v1 := &version{"", 0, 0, 0}
	v2 := &version{"", 1, 0, 0}
	c.Check(v2.GreaterThan(v1), Equals, true)
	c.Check(v1.GreaterThan(v2), Equals, false)

	v1 = &version{"", 0, 0, 0}
	v2 = &version{"", 0, 1, 0}
	c.Check(v2.GreaterThan(v1), Equals, true)
	c.Check(v1.GreaterThan(v2), Equals, false)

	v1 = &version{"", 0, 0, 0}
	v2 = &version{"", 0, 0, 1}
	c.Check(v2.GreaterThan(v1), Equals, true)
	c.Check(v1.GreaterThan(v2), Equals, false)

	v1 = &version{"", 0, 0, 2}
	v2 = &version{"", 0, 1, 0}
	c.Check(v2.GreaterThan(v1), Equals, true)
	c.Check(v1.GreaterThan(v2), Equals, false)

	v1 = new(version)
	v1.parse("05.03.01")
	v2 = &version{"", 0, 1, 0}
	c.Check(v2.GreaterThan(v1), Equals, false)
	c.Check(v1.GreaterThan(v2), Equals, true)

}

//===================================

func fromStatic(url string) (data []byte, err error) {
	json := `[{"ref": "refs/tags/1.0.0" },{"ref": "refs/tags/2.0.0" }]`
	return []byte(json), nil
}
