// Package timelock implements a time based cache for dependencies
package timelock

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vube/depman/colors"
	"github.com/vube/depman/dep"
	"github.com/vube/depman/util"
)

const (
	// the filename of the cache json
	cacheFileName = ".depman.cache"

	// timeout in hours
	timeoutHours = 1.0
)

// flags
var (
	clear bool
	skip  bool
)

var (
	// the cache, a map of repos urls to timestamps
	cache map[string]time.Time

	// the full path to the cache json
	cacheFile string
)

func init() {
	flag.BoolVar(&clear, "clear-cache", false, "Delete the time based cache")
	flag.BoolVar(&skip, "skip-cache", false, "Skip the time based cache for this run only")
}

// Clear clears the cache, returns true if the cache was cleared
func Clear() (cleared bool) {
	cleared = clear

	if clear {
		parts := strings.Split(os.Getenv("GOPATH"), ":")
		cacheFile = filepath.Join(parts[0], cacheFileName)

		util.Print(colors.Yellow("Clearing cache file: " + cacheFile))

		_, err := os.Stat(cacheFile)
		if err != nil {
			return
		}

		err = os.Remove(cacheFile)
		if err != nil {
			util.Fatal(err)
		}
	}

	return
}

// Read reads the cache from disk
func Read() {
	parts := strings.Split(os.Getenv("GOPATH"), ":")
	cacheFile = filepath.Join(parts[0], cacheFileName)

	cache = make(map[string]time.Time)

	if !util.Exists(cacheFile) {
		return
	}

	util.Verbose("Reading cache file from " + cacheFile)

	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		util.Fatal(err)
	}

	err = json.Unmarshal(data, &cache)
	if err != nil {
		return
	}
}

// Write writes the cache out to disk
func Write() {
	if skip {
		return
	}

	var buf bytes.Buffer
	str, err := json.Marshal(cache)
	if err == nil {
		json.Indent(&buf, str, "", "    ")
		data := []byte(buf.String() + "\n")

		util.Verbose("Writing cache file to " + cacheFile)

		err = ioutil.WriteFile(cacheFile, data, 0644)
		if err != nil {
			util.Fatal(err)
		}
	}
	return
}

// IsStale returns true if the cached dependency is older than timeoutHours
func IsStale(d *dep.Dependency) (stale bool) {

	if skip || d.SkipCache {
		return true
	}

	ts, ok := cache[d.Repo]

	// item is in the cache
	if ok {
		// item is old
		if time.Since(ts).Hours() > timeoutHours {
			stale = true
			cache[d.Repo] = time.Now()
		}
	} else {
		stale = true
		cache[d.Repo] = time.Now()
	}

	return
}
