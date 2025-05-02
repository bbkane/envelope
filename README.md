# enventory

![./demo.gif](./demo.gif)

Store environment variables for projects in a central SQLite database!

- Automatically export/unexport environments when entering/leaving directories
- Need an environment variable in more than one environment? Create a reference to it instead of copying it.
- Currently only supports `zsh`

## Project Status

I'm using `enventory` personally, and it seems to work well! That said, I work
on `enventory` for fun, and part of that is changing APIs and CLI interfaces
when I want to.

## Install

- [Homebrew](https://brew.sh/): `brew install bbkane/tap/enventory`
- [Scoop](https://scoop.sh/):

```
scoop bucket add bbkane https://github.com/bbkane/scoop-bucket
scoop install bbkane/enventory
```

- Download Mac/Linux/Windows executable: [GitHub releases](https://github.com/bbkane/enventory/releases)
- Go: `go install go.bbkane.com/enventory@latest`
- Build with [goreleaser](https://goreleaser.com/) after cloning: ` goreleaser release --snapshot --clean`

## Initialize in `~/.zshrc`

> Other shells not yet supported

```bash
eval "$(enventory shell zsh init)"
```

## Initialize `zsh` Tab Completion

`enventory` is quite a verbose CLI, so tab completion (which also auto-completes env names, var names, and var ref names) is super useful.

```bash
enventory completion zsh > /something/in/$fpath
```

## Dev Notes

See [Go Project Notes](https://www.bbkane.com/blog/go-project-notes/) for notes on development tooling and CI/CD setup (including demo gif generation)

### Generate [`./dbdoc`](./dbdoc) with [tbls](https://github.com/k1LoW/tbls)

Install:

```bash
brew install k1LoW/tap/tbls
```

Run:

```bash
# get a fresh db
go run . env list --db-path tmp.db
tbls doc --rm-dist
```

### Generate [./sqlite/sqlite/sqlcgen](./sqlite/sqlite/sqlcgen)

```bash
go generate ./...
```

