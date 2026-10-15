# pg-query-watch

Show PostgreSQL statements with high mean execution time from the pg_stat_statements snapshot.

## Install

Requires Go 1.23 or newer and PostgreSQL client tools. Build with `go install github.com/ulrichheringer/pg-query-watch@latest`.

## Use

Enable pg_stat_statements, configure standard PGHOST/PGDATABASE/PGUSER variables, then run `pg-query-watch -min-ms 100`. Set PGPASSWORD through your secret manager or shell environment.

## Example

Run `go run . -limit 20 -min-ms 50` against a PostgreSQL database with pg_stat_statements enabled.
