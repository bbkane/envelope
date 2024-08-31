# Changelog

All notable changes to this project will be documented in this file. The format
is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

Note the the latest version is usually work in progress and may have not yet been released.

# v0.0.10

## Changed

- Skip printing Comment if it's blank
- Skip printing UpdateTime if it equals CreateTime

# v0.0.9

## Added

- `init zsh` has a `--print-autoload` flag now
- `env var update` command

## Changed

- `env var create --value` is optional, and the value is prompted for if not given 
- `init` -> `init zsh` so we can add zsh-specific flags and make subcommands for other shells

# v0.0.8

## Added

- `--format` flag to change output format (currently only supports the default (`table` and `value-only` for vars and refs))

## Fixed

- `--mask` flag now hides values in `env var show`

# v0.0.7

## Added

- `--mask` flag to show commands to hide sensitive values

# v0.0.6

## Added

- `init` now takes flags to gate stuff to print
- `init` now adds `export-env` and `unexport-env` to the environment

## Fixed

- Fixed spelling for `env ref create`
- Unexport `$OLDPWD` env before exporting `$PWD`, so if they share an export name, the new one isn't deleted

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
