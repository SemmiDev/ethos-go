// Package migrations provides embedded SQL migration files for golang-migrate.
// This package must be at the project root level to embed the migrations folder.
package migrations

import "embed"

// FS contains all SQL migration files embedded at build time.
// Usage: database.RunMigrations(dbURL, migrations.FS, ".")
//
//go:embed *.sql
var FS embed.FS
