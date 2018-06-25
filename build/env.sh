#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
lbchain-devdir="$workspace/src/github.com/lbchain-devereum"
if [ ! -L "$lbchain-devdir/go-lbchain-devereum" ]; then
    mkdir -p "$lbchain-devdir"
    cd "$lbchain-devdir"
    ln -s ../../../../../. go-lbchain-devereum
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$lbchain-devdir/go-lbchain-devereum"
PWD="$lbchain-devdir/go-lbchain-devereum"

# Launch the arguments with the configured environment.
exec "$@"
