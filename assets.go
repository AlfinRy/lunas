// Package assets embeds the built SPA and hands it to the server at deploy
// time, so the whole product ships as a single binary.
package assets

import (
	"embed"
	"io/fs"
)

//go:embed all:web/dist
var webFS embed.FS

//go:embed db/schema.sql
var schemaSQL string

// SchemaSQL exposes the embedded schema for main.
func SchemaSQL() string { return schemaSQL }

// Web returns the SPA files rooted at web/dist.
func Web() fs.FS {
	sub, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		panic(err) // compile-time embed; unreachable
	}
	return sub
}
