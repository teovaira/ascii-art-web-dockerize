// Package banners provides the embedded ASCII art banner files for the web server.
//
// The three banner files (standard, shadow, thinkertoy) are compiled into
// the binary at build time using go:embed. This makes the binary fully
// self-contained and relocatable — it can be run from any directory without
// requiring the banner files to exist on disk at runtime.
//
// Usage:
//
//	bannerData, err := parser.LoadBanner(banners.FS, "standard.txt")
package banners

import "embed"

// FS is the embedded filesystem containing all banner .txt files.
//
// Files included at compile time: standard.txt, shadow.txt, thinkertoy.txt
//
//go:embed *.txt
var FS embed.FS
