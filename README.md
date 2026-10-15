# PostgreSQL Query Watch

[![Go](https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/) [![MIT License](https://img.shields.io/badge/License-MIT-007ec6.svg)](LICENSE)

Find PostgreSQL statements with high average execution time using pg_stat_statements.

## Install

Requires Go 1.23 or newer and the PostgreSQL client tools. Install with `go install github.com/ulrichheringer/pg-query-watch@latest`.

## Use

Enable pg_stat_statements and set `PGHOST`, `PGDATABASE`, and `PGUSER` in your environment. Add `PGPASSWORD` if needed, then run `pg-query-watch -min-ms 100`.

## Example

Run `go run . -limit 20 -min-ms 50` against a database with pg_stat_statements enabled.
