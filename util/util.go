// Package util provides various utility functions.
// It also defines the following flags: --verbose,  --debug, --silent, and --version
package util

// Copyright 2013 Vubeology, Inc.

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	pathPkg "path"
	"strings"

	"github.com/vube/depman/colors"
	"github.com/vube/depman/dep"
)

// GLOBALS
var (

	// Whether to display the version number and exit
	displayVersion bool

	// Whether to display debug messages
	debug bool

	// Whether to display verbose messages
	verbose bool

	// don't display any output
	silent bool

	// display help message
	help bool

	// current indent level
	indentLevel int

	// Sub command like install add or update
	command string
)

//===============================================

// Polymorphics to allow mocking in tests, these should always be set to their defaults except during testing
var (
	Fatal      func(v ...interface{})
	RunCommand = defaultRun
	OsExit     = os.Exit
	Cd         = defaultCd
	indent     = defaultIndent
)

var (
	// whether to parse flags or not, this should always be set to true except during tests
	parse = true

	logger *log.Logger

	// OutputTarget is an io.Writer to write messages to, change it to a bytes.buffer to test output
	OutputTarget = os.Stderr
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Display debug messages. Implies --verbose")
	flag.BoolVar(&verbose, "verbose", false, "Display commands as they are run, and other informative messages")
	flag.BoolVar(&silent, "silent", false, "Don't display normal output. Overrides --debug and --verbose")
	flag.BoolVar(&displayVersion, "version", false, "Display version number")
	logger = log.New(OutputTarget, "", 0)
	Fatal = logger.Fatal
}

// Parse command line flags
func Parse() {
	flag.Parse()

	if silent {
		debug = false
		verbose = false
	}

	if debug {
		verbose = true
		logger.SetFlags(log.Lshortfile)
	}
}

//GetPath processes p and returns a clean path ending in deps.json
func GetPath(p string) (path string) {
	if !strings.HasSuffix(p, dep.DepsFile) {
		path = p + "/" + dep.DepsFile
	}
	path = pathPkg.Clean(path)
	return
}

// Version displays the version of depman and optionally exits (--version)
func Version(v string) {
	if displayVersion {
		logger.Output(2, "Depman Version "+v)
		OsExit(0)
	} else if !silent {
		logger.Output(2, "Depman Version "+v)
	}
}

// Change directory to the specified directory, checking for errors
// Returns the path to the old working directory
func defaultCd(dir string) (result int) {
	err := os.Chdir(dir)

	if err != nil {
		logger.Output(2, indent()+colors.Red("$ cd "+dir))
		logger.Output(2, indent()+colors.Red(err.Error()))
		result = 1
	} else if verbose {
		logger.Output(2, indent()+"$ cd "+dir)
	}
	return
}

// Pwd returns the current working directory
func Pwd() (pwd string) {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Print(colors.Red("Cannot get Current Working Directory"))
		Fatal(colors.Red(err.Error()))
	}
	return
}

// Wrapper on os.exec to catch errors, and print useful messages
func defaultRun(cmd string) (result int) {

	if verbose {
		logger.Output(2, indent()+"$ "+cmd)
	}

	parts := strings.Split(cmd, " ")
	c := exec.Command(parts[0], parts[1:]...)

	out, err := c.CombinedOutput()

	if err != nil {
		logger.Output(2, indent()+colors.Red("$ "+cmd))
		logger.Output(2, indent()+colors.Red(string(out)))
		logger.Output(2, indent()+colors.Red(err.Error()))
		result += 1
	}

	if len(out) > 0 && debug {
		logger.Output(2, indent()+string(out))
	}
	return
}

// UpwardFind searches for file starting in dir and moving up the path.
// Returns the path to the found file or the empty string if it was not found
func UpwardFind(dir string, file string) (found string) {
	split := strings.Split(dir, "/")
	for skip := len(split); skip >= 0; skip-- {
		f := strings.Join(split[:skip], "/") + "/" + file
		if Exists(f) {
			found = f
			return
		}
	}
	return
}

// Exists checks if a file or directory exists
func Exists(path string) (res bool) {
	_, err := os.Stat(path)
	if err == nil {
		res = true
	}
	return
}

// IncreaseIndent Increments the indentation level used during PrintIndent calls
func IncreaseIndent() {
	indentLevel++
}

// DecreaseIndent decrements the indentation level used during PrintIndent calls
func DecreaseIndent() {
	indentLevel--
}

func defaultIndent() (res string) {
	res = strings.Repeat(" |", indentLevel+1) + " "
	return
}

// PrintDep displays a dependency based on the --silent and --verbose flags
func PrintDep(name string, d dep.Dependency) {
	if !silent {
		if verbose {
			logger.Output(2, indent()+colors.Blue(name)+colors.Yellow(" ("+d.Version+")")+" "+d.Repo)
		} else {
			logger.Output(2, indent()+colors.Blue(name)+colors.Yellow(" ("+d.Version+")"))
		}
	}
}

// CheckPath causes the application to exit if path does not exist
func CheckPath(path string) {
	if !Exists(path) {
		Fatal(colors.Red("Could not find '" + path + "' are you in the right directory?"))
	}
}

// Mock configures flags indentation and the logger for testing
func Mock(w io.Writer) {
	verbose = false
	debug = false
	silent = false
	logger = log.New(w, "", 0)
	Fatal = logger.Print
	indent = func() string {
		return ""
	}

}

// Verbose prints s if --verbose is set
func Verbose(s string) {
	if verbose {
		logger.Output(2, s)
	}
}

// VerboseIndent prints s with indentation if --verbose is set
func VerboseIndent(s string) {
	if verbose {
		logger.Output(2, indent()+s)
	}
}

// Print prints s unless --silent is set
func Print(s string) {
	if !silent {
		logger.Output(2, s)
	}
}

// PrintIndent prints s with indentation unless --silent is set
func PrintIndent(s string) {
	if !silent {
		logger.Output(2, indent()+s)
	}
}

// Debug prints s id --debug is set
func Debug(s string) {
	if debug {
		logger.Output(2, s)
	}
}

// SetVerbose sets the verbose level for testing
func SetVerbose(v bool) {
	verbose = v
}
