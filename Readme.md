# depman
--
Dependency management helper for Golang packages. Supports versioned
dependencies using standard Golang imports.


### Installation

Run:

    go get github.com/vube/depman
    depman help


### Author

Nicholas Capo <nicholas@vubeology.com>


### Copyright

Copyright 2013-2014 Vubeology, Inc. Released under the MIT License (see
LICENSE).


### Features

* Simple json configuration file

* Recursive installs

* Uses `go get` when possible

* Supports Git, Mercurial, and Bazaar

* Automatic installation of dependencies from git repositories that do not
support `go get`

* Time based "cache" to increase speed

* Handles Multi-Part $GOPATH


### Usage

Run `depman` in the directory with your `deps.json`, or use the `--path`
argument to specify the directory where `deps.json` is located.


### Commands

* `init` Create an empty deps.json

* `add` Add a dependency (interactive)

* `install` Install all the dependencies listed in deps.json (default)

* `update [nickname] [branch]` Update [nickname] to use the latest commit in
[branch]

* `show-frozen` Show dependencies as resolved to commit IDs. Use the
`--recursive` flag to descend into dependencies depth-first.

* `help` Display help message


### Options

* `-clean=false`: Remove changes to code in dependencies

* `-clear-cache=false`: Delete the time based cache

* `-debug=false`: Display debug messages. Implies --verbose

* `-help=false`: Display help

* `-no-colors=false`: Disable colors

* `-path="."`: Directory or full path to deps.json

* `-silent=false`: Don't display normal output. Overrides --debug and --verbose

* `-skip-cache=false`: Skip the time based cache for this run only

* `-verbose=false`: Display commands as they are run, and other informative
### messages

* `-version=false`: Display version number


### Basic Algorithm

For each dependency:

1. Use `go get` to download the dependency

2. Checkout the specified version

3. Look in the dependency's directory for a `deps.json` and recursively install
those dependencies

4. If the dependency type is `git-clone` then manually run `git clone`, `git
fetch`, etc as needed


### Duplicates

Duplicated dependencies (anywhere in the tree) with _identical_ versions will be
skipped (this just saves some time and prevents infinite recursion). Duplicated
dependencies with _different_ versions cause a fatal error and must be fixed by
the developer.


### Non Go-Getable Repos

Some repositories (private bitbucket repositories for example), are not
supported by `go get`. To include those repositories in depman:

1. Change the type to `git-clone` (hg and bzr are not yet supported)

2. Change the `repo` to a full git url (include everything necessary for `git
clone`)

3. Add an `alias` field to specify a directory in which to clone, the path is
rooted at `$GOPATH/src/`

See the example below or the included `deps.json` file.


Multi-Part $GOPATH Support

Depman since version 2.8.0 supports multi-part $GOPATH. When installing
dependencies, depman will search the parts of $GOPATH for the appropriate
directory, if not found, depman will install the dependency to the first part of
$GOPATH.


### Cache

Depman uses a time based cache to speed up normal operations.

A global list of dependencies and timestamps is kept at `$GOPATH/.depman.cache`.
When depman is run it looks at the timestamp to decide whether to update the
repo or not (`go get -u`, `git clone`, `git fetch`, etc).

If the dependency is more that 1 hour old, depman will fetch updates from the
network, otherwise depman uses the repo as is.

If the cache was stale or unused, a '*' will be printed at the end of the
installation line.

You can clear the cache by deleting the cache file, or running depman with the
`--clear-cache` flag. The cache can be skipped for the current run by using the
`--skip-cache` flag.

Additional information about the cache (including the time spent while
installing) can been seen by running depman with the `--verbose` flag.

The code for this feature is in the `timelock` package.


### Implementation Requirements

* Depman shall not require any external dependencies (beyond the standard
library) for normal operation

* Testing dependencies shall be installed through depman

* `go vet` shall not indicate any problems

* `golint` (https://github.com/golang/lint) shall not indicate any problems


### Testing

Depman can be tested by running `make` in the source directory.

This does the following:

1. Installs depman

2. Installs testing dependencies

3. Runs depman with a few common flags to confirm basic operation

4. Run `go vet`, and `golint`

5. Executes the unit tests

6. Generate `Readme.md`

See `Makefile` for more information.


### JSON Structure

NOTE: The file must be named `deps.json`

    {
    	"shortname":{
    		"repo":"url/to/package, just like in import",
    		"version":"commit, tag, or branch",
    		"type": "one of 'git', 'bzr', 'hg'"
    		"skip-cache":"optional, set to 'true' to always ignore the cache"
    	},
    	"not go getable":{
    		"repo":"full git repo url, just like git clone"
    		"version":"commit, tag, or branch",
    		"type": "must be 'git-clone'",
    		"alias": "target directory to clone into, (only supported for type 'git-clone')"
    	}
    }

### Example

    {
    	"gocheck": {
    		"repo": "launchpad.net/gocheck",
    		"version": "87",
    		"type": "bzr"
    	},
    	"gocov": {
    		"repo": "github.com/vube/gocov/gocov",
    		"version": "93371a7ae85bec1c4afe9b9f3281c062ab106e6d",
    		"type": "git",
    		"skip-cache": true
    	},
    	"gocov-html": {
    		"repo": "https://github.com/matm/gocov-html.git",
    		"version": "1512341d22ab06788fc5ad63925fd07979a9ef39",
    		"type": "git-clone",
    		"alias": "github.com/matm/gocov-html"
    	},
    	"godocdown": {
    		"repo": "github.com/robertkrimen/godocdown",
    		"version": "0bfa0490548148882a54c15fbc52a621a9f50cbe",
    		"type": "git"
    	},
    	"golint": {
    		"repo": "github.com/golang/lint/golint",
    		"version": "7e9cdc93310598b5c6cd9ce4d0d0a37c0f5b9e4c",
    		"type": "git"
    	}
    }


### Todo

* Support manually cloning `hg` or `bzr` repositories (add types `hg-clone`, and
`bzr-clone`)
