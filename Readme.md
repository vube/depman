Depman
=====

Dependency management helper for Golang packages. Supports versioned dependencies, using standard Golang imports.

Copyright 2013 Vubeology, Inc.

Installation
----------------

	git clone git@github.com:vube/depman.git $GOPATH/src/vube/depman
	go install github.com/vube/depman

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
* `--clean=false`: Remove changes to code in dependencies
* `--debug=false`: Display debug messages. Implies --verbose
* `--help=false`: Display help
* `--no-colors`=false: Disable colors
* `--path="./":` Directory or full path to .deps.json
* `--silent=false`: Don't display normal output. Overrides --debug and --verbose
* `--verbose=false`: Display commands as they are run, and other informative messages
* `--version=false`: Display version number

Requirements
---------------------

* Shall not require any external dependencies for normal operation
* Testing dependencies shall be installed through depman
* `go vet` shall not indicate any problems
* `golint` (github.com/golang/lint) shall not indicate any problems

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

	# Lint
	golint .

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

JSON Structure
----------------------

[The file must be named `deps.json`]

	{
		"shortname":{
			"repo":"url/to/package, just like in import",
			"version":"commit, tag, or branch",
			"type": "one of 'git', 'bzr', 'hg'"
		},
		"not go getable":{
			"repo":"full git repo url, just like git clone"
			"version":"commit, tag, or branch",
			"type": "git-clone",
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
			"repo": "github.com/axw/gocov/gocov",
			"version": "90c9390",
			"type": "git"
		},
		"gocov-html": {
			"repo": "https://github.com/matm/gocov-html.git",
			"version": "f83448b",
			"type": "git-clone",
			"alias": "github.com/matm/gocov-html"
		},
		"golint": {
			"repo": "github.com/golang/lint",
			"version": "22f4d5e",
			"type": "git"
		}
	}

TODO
--------

* Support manually cloning `hg` or `bzr` repositories (add types `hg-clone`, and `bzr-clone`)
