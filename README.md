# envelope

![./demo.gif](./demo.gif)

Store environment variables for projects in a central SQLite database!

- Automatically export/unexport environments when entering/leaving directories
- Need an environment variable in more than one environment? Create a reference to it instead of copying it.

## Project Status

I'm using `envelope` personally, but I can't recommend it for anyone else to use until I have more features and tab completion. The CLI interface is also not stable.

## Install

- [Homebrew](https://brew.sh/): `brew install bbkane/tap/envelope`
- [Scoop](https://scoop.sh/):

```
scoop bucket add bbkane https://github.com/bbkane/scoop-bucket
scoop install bbkane/envelope
```

- Download Mac/Linux/Windows executable: [GitHub releases](https://github.com/bbkane/envelope/releases)
- Go: `go install go.bbkane.com/envelope@latest`
- Build with [goreleaser](https://goreleaser.com/) after cloning: `goreleaser --snapshot --skip-publish --clean`

## Initialize in `~/.zshrc`

> Other shells not yet supported

```bash
eval "$(envelope shell zsh init)"
```

## Dev Notes

### Generate [./sqlite/sqlite/sqlcgen](./sqlite/sqlite/sqlcgen)

```bash
go generate ./...
```

See [Go Project Notes](https://www.bbkane.com/blog/go-project-notes/) for notes on development tooling and CI/CD setup.

### Generate [`./dbdoc`](./dbdoc) with [tbls](https://github.com/k1LoW/tbls)

Install:

```bash
brew install k1LoW/tap/tbls
```

Run:

```bash
# get a fresh db
envelope env list --db-path tmp.db
tbls doc --rm-dist
```

