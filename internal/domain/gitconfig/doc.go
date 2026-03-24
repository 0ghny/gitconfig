// Package gitconfig manages read and write operations on a single .gitconfig
// file. It provides the GitConfig value object backed by an injected afero
// filesystem, making it usable with both the real OS filesystem and an
// in-memory filesystem for testing.
package gitconfig
