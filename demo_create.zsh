#!/usr/bin/env zsh

# exit the script on command errors or unset variables
# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

# https://stackoverflow.com/a/246128/295807
script_dir="${0:A:h}"
readonly script_dir
cd "${script_dir}"

# Use a new version of envelope
go install .
export PATH="/Users/bbkane/go/bin:$PATH"

export PROMPT='%F{47}$ %f'
rm -f ./tmp.db
export ENVELOPE_DB_PATH=./tmp.db

vhs < ./demo.tape
