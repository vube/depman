package timelock

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

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

func Read() {
	cacheFile = filepath.Join(os.Getenv("GOPATH"), cacheFileName)

	cache = make(map[string]time.Time)

	if !util.Exists(cacheFile) {
		return
	}

	if clear {
		err := os.Remove(cacheFile)
		if err != nil {
			util.Fatal(err)
		}
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

func IsStale(repo string) (stale bool) {

	if skip {
		return true
	}

	ts, ok := cache[repo]

	// item is in the cache
	if ok {
		// item is old
		if time.Since(ts).Hours() > timeoutHours {
			stale = true
			cache[repo] = time.Now()
		}
	} else {
		stale = true
		cache[repo] = time.Now()
	}

	return
}
