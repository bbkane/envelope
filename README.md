# envelope

![./demo.gif](./demo.gif)

Store environment variables for projects relationally! UX inspired by the Azure CLI.

Expand details below to see the planned UX

<details>

Not all of these commands are implemented yet (`go run . -h outline` to see what is), but, in general, the UX will look like:

```bash
envelope
    --database ~/.config/envelope.db
    env
        create
            --name $PWD
            --comment "blah"
        delete
            --name $PWD
        list
        print-script
            --name $PWD
            --include-parent-dirs true # try to find an env in parent dirs
            --shell zsh
            --type export | unexport
        show --name $PWD  # Also shows all vars that will be exported
        update ...
        keyringref
            --env-name $PWD
            create
                --env-ref-name other_env_name
                --localvar-ref-name envvar_name
            delete ...
            list
                --env-name
            show
            update ...
        localvar
            --env-name
            create
                --name bob
                --value bob
            delete ...
            list
                --env-name
            show
            update ...
        localref
            --env-name $PWD
            create
                --env-ref-name other_env_name
                --localvar-ref-name envvar_name
            delete ...
            list
                --env-name
            show
            update ...
    keyring
        create --name azure_client_secret ... # prompt for value
        show --name
        update --name ...
        delete --name azure_client_secret
        list --print-values true
```

</details>

# Questions for Manuel

Thanks for the help!!

## Is the amount of code for this normal?

I've got a "layered" architecture:

- CLI layer that reads os.Args and calls the "domain" layer, then prints it nicely
- "domain" layer that's basically an interface (maybe I should call it a "storage" or "persistence" layer)
- "sqlite" layer that implements the "domain" layer interface

I've got gobs of code to translate between all of these layers and I'm like 40% done with commands I need to implement.

The process I have to add a new command is:

- Update these notes to get a CLI I like
- Update SQL and generate code with [`sqlc`](https://sqlc.dev/).
- Update `sqlite.EnvService` to use generated code.
- Update `domain.EnvService` with the new method added to `sqlte.EnvService`
- Update `cli` to call `domain.EnvService`
- Update `main` to add the new command to the app tree.
- Update `main_test_<table>.go` to ensure everything works (this generates testdata when `ENVELOPE_TEST_UPDATE_GOLDEN=1 go test ./...` is run)

Is there a better way to organize this than what I have? Or is this fairly normal?

## My [`EnvService`](./domain/env.go) interface is getting huge

In theory, I'd like to have this much smaller, but in practice, all operations need to touch the same db, and all the tables are connected, and this code gatekeeps access to that db, so maybe it needs a all the methods?

Is there a better way to organize this?

## What do you think of the idea generally?

I think storing env vars in a SQLite DB is a good idea, but I won't *really* know until this is functional enough to use (i.e., keyring support finished and hooked up to environments).

Then I can hook this into `zsh`'s [chpwd](https://stackoverflow.com/a/3964198/2958070) function to automatically add envvars when I enter a directory and remove them when I leave the directory. Similar to `direnv`, but centralized management.

## Any other comments?

Thanks again!!

---

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

## Notes

See [Go Project Notes](https://www.bbkane.com/blog/go-project-notes/) for notes on development tooling.
