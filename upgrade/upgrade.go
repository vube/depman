// Package upgrade provides functions to check for and upgrade depman its self
// Check() and Print() provides a way to asynchronously checks for new versions on github
// Self() provides a way to upgrade depman using depman
package upgrade

// Copyright 2013-2014 Vubeology, Inc.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/vube/depman/colors"
	"github.com/vube/depman/dep"
	"github.com/vube/depman/result"
	"github.com/vube/depman/timelock"
	"github.com/vube/depman/util"
)

const (
	message = "Depman Version %s is available!\nRun 'depman self-upgrade' to upgrade."
	apiURL  = "https://api.github.com/repos/vube/depman/git/refs/tags"
	none    = ""
)

var (
	checkCalled bool
	selfCalled  bool
	channel     chan string
	checkError  error

	self *dep.Dependency
)

func init() {

}

// Self upgrades this version of depman to the latest on the master branch
func Self(version string) {
	selfCalled = true
	util.Print(colors.Blue("Upgrading depman..."))
	util.RunCommand("go get -u github.com/vube/depman")

	cmd := exec.Command("depman", "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		result.RegisterError()
		util.Print(colors.Red(string(out)))
		return
	}

	newVersion := strings.TrimSuffix(strings.TrimPrefix(string(out), "Depman Version "), "\n")

	if newVersion != version {
		util.Print("Upgraded to Version " + newVersion)
	} else {
		util.Print("No upgrade found")
	}
}

//============================================

// Check checks for a newer version of depman on github, the result is written on a channel
// It is expected that this will be called in a goroutine
func Check(ver string) {
	var str string
	checkCalled = true
	channel = make(chan string, 0)

	self = new(dep.Dependency)
	self.Repo = "depman internal upgrade check"

	if timelock.IsStale(self) {
		str, checkError = check(ver)
	} else {
		str = none
	}

	channel <- str
}

// check() does the actual work of checking for a new version
// it is pulled out of Check() to ease testing
func check(ver string) (result string, err error) {
	var current = new(version)

	err = current.parse(ver)
	if err != nil {
		return
	}

	versions, err := getVersions()
	if err != nil {
		return
	}

	m := max(versions)

	if m.GreaterThan(current) {
		result = m.str
	}

	return
}

// Print displays a message if a newer version of depman is available
// this should only be called after Check()
func Print() {
	if !checkCalled {
		util.Fatal(colors.Red("You must call upgrade.Check() before upgrade.Print()"))
	}

	// if the command used was self-upgrade, we can just return
	if selfCalled {
		return
	}

	ref := <-channel

	if checkError != nil {
		util.Verbose(colors.Yellow("Upgrade Check Error: " + checkError.Error()))
	}

	if ref != none {
		fmt.Println(colors.Yellow(fmt.Sprintf(message, ref)))
	}
}

type version struct {
	str string
	a   int
	b   int
	c   int
}

// GreaterThan returns (v > o)
func (v *version) GreaterThan(o *version) bool {
	if v.a > o.a {
		return true
	} else if v.a < o.a {
		return false
	}

	if v.b > o.b {
		return true
	} else if v.b < o.b {
		return false
	}

	if v.c > o.c {
		return true
	} else if v.c < o.c {
		return false
	}

	return false
}

// max returns the greatest version
func max(vers []*version) (m *version) {
	m = new(version)
	for _, v := range vers {
		if v.GreaterThan(m) {
			m = v
		}
	}
	return
}

// parse parses a version number (three ints) from a string
func (v *version) parse(ver string) (err error) {
	v.str = ver
	parts := strings.Split(ver, ".")

	if len(parts) != 3 {
		return fmt.Errorf("versions must be three ints separated by '.'")
	}
	_, err = fmt.Sscanf(ver, "%d.%d.%d", &v.a, &v.b, &v.c)

	return
}

// getVersions fetches ref data using getter() unmarshalls it and returns a slice of versions
func getVersions() (vers []*version, err error) {
	type ref struct {
		Tag string `json:"ref"`
	}

	var refs []ref

	data, err := getter(apiURL)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &refs)
	if err != nil {
		return
	}

	for _, v := range refs {
		p := strings.Split(v.Tag, "/")
		r := p[len(p)-1]

		v := new(version)
		err = v.parse(r)
		if err != nil {
			continue
		}

		vers = append(vers, v)
	}

	return
}

// getter abstracts fetching from github, to make testing easier
// this implementation fetches version data from the gihub api, using a 500 millisecond timout
var getter = func(url string) (data []byte, err error) {
	var timeout = time.Duration(500 * time.Millisecond)

	var transport = http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
	}

	var client = http.Client{
		Transport: &transport,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	return
}
