# envelope

![./demo.gif](./demo.gif)

Store environment variables for projects relationally! UX inspired by the Azure CLI.

## Project Status

I'm using `envelope` personally, but I can't recommend it for anyone else to use until I have more features (update commands, keyringref support) and tab complete/TUI.

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
        var
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
eval "$(envelope init zsh)"
```

## Notes

See [Go Project Notes](https://www.bbkane.com/blog/go-project-notes/) for notes on development tooling.
