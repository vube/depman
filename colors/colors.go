// Package colors provides functions to optionally wrap strings with ASCII colors codes.
// The display of colors is disabled by the --no-colors flag
package colors

// Copyright 2013 Vubeology, Inc.

import (
	"flag"
	"fmt"
)

var (
	// Disable colors
	noColors bool
)

func init() {
	flag.BoolVar(&noColors, "no-colors", false, "Disable colors")
}

// Yellow returns s wrapped in Yellow ASCII Color Codes
func Yellow(s string) (res string) {
	return color(s, "\033[33m")
}

// Red returns s wrapped in Red ASCII Color Codes
func Red(s string) (res string) {
	return color(s, "\033[31m")
}

// Blue returns s wrapped in Blue ASCII Color Codes
func Blue(s string) (res string) {
	return color(s, "\033[36m")
}

// colors returns s prepended with the color code c, and appended with the end color code
func color(s string, c string) (res string) {
	if noColors {
		res = s
	} else {
		res = fmt.Sprintf("%s%s\033[0m", c, s)
	}
	return
}

// Mock disables colors for testing
func Mock() {
	noColors = true
}
