# Changelog

All notable changes to this project will be documented in this file. The format
is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

Note the the latest version is usually work in progress and may have not yet been released.

# v0.0.5

## Added

- `env ref` commands
- `print-script --shell` flag

## Changed

- `env localvar` commands renamed to `env var`
- Use key-value tables for output
- Show `env ref`s in `env show`
- Export `env ref`s in `env export`
- Show `env ref`s in `env var show`
- When listing the same type of item, print a single table with multiple sections instead of separate tables

# v0.0.4

## Added

- `--confirm` flag to deletes / updates
- `envelope env print-script --type unexport`
- `envelope init`

## Changed

- `--sqlite-dsn` -> `--db-path`. Reads from `ENVELOPE_DB_PATH` env var now too
- made all tests parallel
- more concise date format
- use `--help detailed` by default

# v0.0.3

## Added

- `--no-env-no-problem` flag
