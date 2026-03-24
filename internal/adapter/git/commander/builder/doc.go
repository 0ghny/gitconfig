// Package builder provides a fluent builder for constructing os/exec.Cmd
// instances with optional arguments, environment variables, and timeouts.
// It always inherits the OS environment when custom variables are set,
// ensuring PATH and other critical variables remain available.
package builder
