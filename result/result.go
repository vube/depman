// Package result hold a "global" result boolean
// This allows the application to eventually exit non-zero
// but only after completing a task where a portion of the task returned errors
package result

var err bool

// RegisterError should be called when a non-fatal error occurred, exicution should continue but we want to exit non-zero eventually
func RegisterError() {
	err = true
}

// ShouldExitWithError can be called to determine if the application should exit non-zero
func ShouldExitWithError() bool {
	return err
}
