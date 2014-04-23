Depman
=====

Dependency management helper for Golang packages. Supports versioned dependencies, using standard Golang imports.

Copyright 2013-2014 Vubeology, Inc.

Released under the MIT License (see LICENSE).

Features
------------

* Simple json configuration
* Recursive installs
* Uses `go get` when possible
* Supports Git, Mercurial, or Bazaar
* Automatic installation of dependencies from git repositories that do not support `go get`
* Time based cacheing to increase speed
* Handles Multi-Part $GOPATH


Installation
----------------

	go get github.com/vube/depman

Usage
---------

Run `depman` in the directory with your `deps.json`, or use the `--path` argument to specify the directory where `deps.json` is located.

### Commands

#### init
Create an empty deps.json

#### add
Add a dependency (interactive)

#### install
Install all the dependencies listed in deps.json (default)

#### update [nickname] [branch]
Update [nickname] to use the latest commit in [branch]

#### show-frozen
Show dependencies as resolved to commit IDs.
##### Options
* `--recursive`: Descend into dependencies depth-first.

#### help
Display help message

### Options
* `-clean=false`: Remove changes to code in dependencies
* `-clear-cache=false`: Delete the time based cache
* `-debug=false`: Display debug messages. Implies --verbose
* `-help=false`: Display help
* `-no-colors=false`: Disable colors
* `-path="."`: Directory or full path to deps.json
* `-silent=false`: Don't display normal output. Overrides --debug and --verbose
* `-skip-cache=false`: Skip the time based cache for this run only
* `-verbose=false`: Display commands as they are run, and other informative messages
* `-version=false`: Display version number

Basic Algorithm
-----------------------

1. For each dependency:
	1. Use `go get` to download the dependency
	2. Checkout the specified version
	3. Look in the dependency's directory for a `deps.json`
		1. Recursively install those dependencies as well
	4. If the dependency type is `git-clone` then manually run `git clone`, `git fetch`, etc as needed

### Duplicates
Duplicated dependencies (anywhere in the tree) with _identical_ versions will be skipped (this just saves some time and prevents infinite recursion). Duplicated dependencies with _different_ versions cause a fatal error and must be fixed by the developer.

### Non Go-Getable Repos
Some repositories (private bitbucket repos for example), are not supported by `go get`.
To include those repositories in depman:

1. Change the type to `git-clone` (hg and bzr are not yet supported)
2. Change the `repo` to a full git url (include everything necessary for `git clone`)
3. Add an `alias` section to specify a directory in which to clone, the path is rooted at `$GOPATH/src/`
See the example below or the included `deps.json` file.

Multi-Part $GOPATH Support
-------------------------------------

Depman since version 2.8.0 supports multi-part $GOPATH. When installing dependencies, depman will search the parts of $GOPATH for the appropriate directory, if not found, depman will install the dependency to the first part of $GOPATH.

Cache
--------

Depman uses a time based cache to speed up normal operations.

A global list of dependencies and timestamps is kept at `$GOPATH/.depman.cache`. When depman is run it looks at the timestamp to decide weather to update the repo (`go get -u`, `git clone`, `git fetch`, etc). If the dependency is more that 1 hour old, depman will fetch updates from the network, otherwise depman uses the repo as is. If the cache was stale or unused, a '*' will be printed at the end of the installation line.

You can clear the cache by deleting the cache file, or running depman with the `--clear-cache` flag. Additional information about the cache (including the time spent while installing) can been seen by running depman with the `--verbose` flag.

The source code for the feature is in `$GOPATH/src/github.com/vube/depman/timecache`

Upgrade Checks
---------------------
Depman will preform an upgrade check of itself, approximatly once every hour. The check is designed to be non-intrusive, and not slow down any normal operations. If a newer version is detected depman will print a helpful message with the new version number and instructions on how to upgrade.

The frequency of checks is controlled by the same cache as described above. Additionally any errors encounters during the check can be viewed by running depman with the `--verbose` flag.

Implementation Requirements
-----------------------------------------

* Depamn shall not require any external dependencies for normal operation
* Testing dependencies shall be installed through depman
* `go vet` shall not indicate any problems
* `golint` (https://github.com/golang/lint) shall not indicate any problems

Testing
-----------

	# Setup
	cd $GOPATH/github.com/vube/depman
	git pull
	go install
	depman install

	# Test
	go test -i ./...
	go test ./...

	# Coverage
	gocov test ./... | gocov report

	# Vet
	go vet ./...

	# Golint
	go get https://github.com/golang/lint
	golint .

JSON Structure
----------------------

 NOTE: The file must be named `deps.json`

	{
		"shortname":{
			"repo":"url/to/package, just like in import",
			"version":"commit, tag, or branch",
			"type": "one of 'git', 'bzr', 'hg'"
			"skip-cache":"optional, set to 'true' to always skip the cache"
		},
		"not go getable":{
			"repo":"full git repo url, just like git clone"
			"version":"commit, tag, or branch",
			"type": "must be 'git-clone'",
			"alias": "target directory to clone into, (optional for every type other than 'git-clone')"
		}
	}

### Example
	{
		"gocheck": {
			"repo": "launchpad.net/gocheck",
			"version": "85",
			"type": "bzr"
			},
		"gocov": {
			"repo": "github.com/vube/gocov/gocov",
			"version": "90c9390",
			"type": "git",
			"skip-cache": true
		},
		"gocov-html": {
			"repo": "https://github.com/matm/gocov-html.git",
			"version": "f83448b",
			"type": "git-clone",
			"alias": "github.com/matm/gocov-html"
		},
		"golint": {
			"repo": "github.com/golang/lint/golint",
			"version": "22f4d5e",
			"type": "git"
		}
	}

Todo
--------

* Support manually cloning `hg` or `bzr` repositories (add types `hg-clone`, and `bzr-clone`)
