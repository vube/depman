package util

// Copyright 2013-2014 Vubeology, Inc.

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/vube/depman/colors"
	. "launchpad.net/gocheck"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

type TestSuite struct{}

var _ = Suite(&TestSuite{})

//===============================================

var PWD string
var GOPATH string

//===============================================

func (s *TestSuite) SetUpSuite(c *C) {
	GOPATH = os.Getenv("GOPATH")
	PWD = Pwd()
	colors.Mock()
}

func (s *TestSuite) SetUpTest(c *C) {
	Cd(PWD)

	debug = false
	verbose = false
	indentLevel = 0
	Cd = defaultCd

	log.SetFlags(0)
	Fatal = log.Println
	log.SetOutput(os.Stdout)
	os.Setenv("GOPATH", GOPATH)
	indent = func() string {
		return ""
	}
}

//===============================================

func (s *TestSuite) TestIndent(c *C) {
	indent = defaultIndent
	indentLevel = 0
	c.Check(indentLevel, Equals, 0)
	c.Check(indent(), Equals, " | ")
	IncreaseIndent()
	c.Check(indentLevel, Equals, 1)
	c.Check(indent(), Equals, " | | ")
	IncreaseIndent()
	c.Check(indentLevel, Equals, 2)
	c.Check(indent(), Equals, " | | | ")
	IncreaseIndent()
	c.Check(indentLevel, Equals, 3)
	c.Check(indent(), Equals, " | | | | ")
	IncreaseIndent()
	c.Check(indentLevel, Equals, 4)
	c.Check(indent(), Equals, " | | | | | ")
	IncreaseIndent()
	c.Check(indentLevel, Equals, 5)
	c.Check(indent(), Equals, " | | | | | | ")

	DecreaseIndent()
	c.Check(indentLevel, Equals, 4)
	DecreaseIndent()
	c.Check(indentLevel, Equals, 3)
	DecreaseIndent()
	c.Check(indentLevel, Equals, 2)
	DecreaseIndent()
	c.Check(indentLevel, Equals, 1)
	DecreaseIndent()
	c.Check(indentLevel, Equals, 0)
}

func (s *TestSuite) TestCdPwd(c *C) {
	Cd("/")
	c.Check(Pwd(), Equals, "/")
	Cd("/etc")
	c.Check(Pwd(), Equals, "/etc")

	verbose = true

	var b []byte
	buf := bytes.NewBuffer(b)
	Mock(buf)
	verbose = true

	err := Cd("/")
	c.Check(buf.String(), Equals, "$ cd /\n")
	c.Check(err, IsNil)

	buf.Truncate(0)
	err = Cd("/none")
	c.Check(buf.String(), Equals, "$ cd /none\nchdir /none: no such file or directory\n")
	c.Check(err, ErrorMatches, `chdir /none: no such file or directory`)
}

func (s *TestSuite) TestPwdErr(c *C) {
	dir, err := ioutil.TempDir("", "DepmanUnitTest")
	c.Check(err, IsNil)
	Cd(dir)
	d := Pwd()
	c.Check(d, Equals, dir)
	os.Remove(dir)

	var b []byte
	buf := bytes.NewBuffer(b)
	Mock(buf)
	d = Pwd()
	c.Check(d, Equals, "")
	c.Check(buf.String(), Equals, "Cannot get Current Working Directory\ngetwd: no such file or directory\n")
}

func (s *TestSuite) TestExists(c *C) {
	c.Check(Exists("/"), Equals, true)
	c.Check(Exists("/etc"), Equals, true)
	c.Check(Exists("./util.go"), Equals, true)
	c.Check(Exists("none"), Equals, false)
	c.Check(Exists("/none"), Equals, false)
}

func (s *TestSuite) TestDefaultRun(c *C) {
	_, err := os.Open("../tests/touch")
	c.Check(err, Not(IsNil))

	RunCommand("touch ../tests/touch")

	_, err = os.Open("../tests/touch")
	c.Check(err, IsNil)

	os.Remove("../tests/touch")

	var b []byte
	buf := bytes.NewBuffer(b)
	Mock(buf)

	verbose = true
	err = defaultRun("echo NNN")
	c.Check(err, IsNil)
	c.Check(buf.String(), Equals, "$ echo NNN\n")

	buf.Truncate(0)
	verbose = false
	debug = true
	err = defaultRun("echo NNN")
	c.Check(err, IsNil)
	c.Check(buf.String(), Equals, "NNN\n")

	buf.Truncate(0)
	verbose = false
	debug = true
	err = defaultRun("none")
	c.Check(buf.String(), Equals, "$ none\n")
	c.Check(err, ErrorMatches, `exec: "none": executable file not found in \$PATH`)
}

func (s *TestSuite) TestUpwardFind(c *C) {
	f := UpwardFind("../tests/success/default", "deps.json")
	c.Check(f, Equals, "../tests/success/default/deps.json")

	f = UpwardFind("../tests/success/default", "main.go")
	c.Check(f, Equals, "../main.go")

	f = UpwardFind("../tests/success/default", "none")
	c.Check(f, Equals, "")
}

func (s *TestSuite) TestPrintDep(c *C) {
	var b []byte
	buf := bytes.NewBuffer(b)
	Mock(buf)

	PrintDep("NNN", "version", "repo", false)
	c.Check(buf.String(), Equals, "NNN (version)\n")

	buf.Truncate(0)
	verbose = true
	PrintDep("NNN", "version", "repo", false)
	c.Check(buf.String(), Equals, "NNN (version) repo\n")

	buf.Truncate(0)
	verbose = true
	PrintDep("NNN", "version", "repo", true)
	c.Check(buf.String(), Equals, "NNN (version) repo *\n")

	buf.Truncate(0)
	verbose = false
	PrintDep("NNN", "version", "repo", true)
	c.Check(buf.String(), Equals, "NNN (version) *\n")
}
