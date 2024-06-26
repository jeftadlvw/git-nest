#!/usr/bin/env bash

SRC_DIR=$1
BIN_DIR=$2
TEST_ENV_DIR=$3

# some default git repositories
export REPO_1="https://github.com/githubtraining/actions-and-packages.git"
export REPO_2="https://github.com/githubtraining/exercise-use-secrets.git"
export REPO_3="https://github.com/githubtraining/hellogitworld.git"
export REPO_4="https://github.com/jeftadlvw/example-repository.git"
export REPO_5="https://github.com/git-lfs/git-lfs.git"

function printUsage() {
    echo "usage: entrypoint.sh [source directory] [binary output directory] [test environment directory]"
}

if [ -z "$SRC_DIR" ]; then
  echo "error: source directory not defined"
  printUsage
  exit 1
fi

if [ -z "$BIN_DIR" ]; then
  echo "error: binary output directory not defined"
  printUsage
  exit 1
fi

if [ -z "$TEST_ENV_DIR" ]; then
  echo "error: test environment directory not defined"
  printUsage
  exit 1
fi

if [ -r "$SRC_DIR" ] && [ -w "$SRC_DIR" ]; then
    echo "error: source directory is not read-only."
    echo "Mount it as readonly to proceed."
    exit 1
fi

function watch_and_build {
    # construct build command
    BUILD_CMD="PWD=$(pwd) && cd $1"

    if [ -z "$DISABLE_TESTS" ]; then
        BUILD_CMD="$BUILD_CMD && go test ./.../tests"
    fi

    BUILD_CMD="$BUILD_CMD && go build -o $2/git-nest -buildvcs=false && cd $PWD"

    echo "Watching directory: $1"
    echo "Output directory: $2"
    echo "Go version: $(go version)"
    echo "Build command: $BUILD_CMD"

    echo "Initial build..."
    eval $BUILD_CMD

    while inotifywait -r -e modify,create,delete,move "$1"; do
        echo "Changes detected. Rebuilding..."
        eval $BUILD_CMD
        echo "=========================================="
    done
}

function prune {
    rm -rf $TEST_ENV_DIR/*
}

function list_repos {
    echo "\$REPO_1=$REPO_1"
    echo "\$REPO_2=$REPO_2"
    echo "\$REPO_3=$REPO_3"
    echo "\$REPO_4=$REPO_4"
    echo "\$REPO_5=$REPO_5"
}

# add binary directory to PATH
PATH="$PATH:$BIN_DIR"

# start tmux session
SESSION_NAME=test-env
tmux new-session -d -s $SESSION_NAME -n "test environment"
tmux set-environment -t $SESSION_NAME -g BIN_DIR $BIN_DIR
tmux set-environment -t $SESSION_NAME -g SRC_DIR $SRC_DIR
tmux set-environment -t $SESSION_NAME -g TEST_ENV_DIR $TEST_ENV_DIR
tmux set-environment -t $SESSION_NAME -g DISABLE_TESTS $TEST_ENV_DIR

# allow mouse input
tmux set -g mouse

# perform horizontal split
tmux split-window -h -t $SESSION_NAME

# rename panes
tmux select-pane -t $SESSION_NAME:0.0 -T "watcher"
tmux select-pane -t $SESSION_NAME:0.1 -T "bash"

# setup path variables in both panes
tmux send-keys -t $SESSION_NAME:0.0 "PATH=$(echo $PATH) && clear && history -c && history -w " C-m
tmux send-keys -t $SESSION_NAME:0.1 "PATH=$(echo $PATH) && clear && history -c && history -w" C-m

# create file watcher
tmux send-keys -t $SESSION_NAME:0.0 "$(declare -f watch_and_build) && clear && history -c && history -w" C-m
tmux send-keys -t $SESSION_NAME:0.0 "watch_and_build $SRC_DIR $BIN_DIR" C-m

# set up prune function in bash pane
tmux send-keys -t $SESSION_NAME:0.1 "$(declare -f prune) && clear && history -c && history -w" C-m
tmux send-keys -t $SESSION_NAME:0.1 "$(declare -f list_repos) && clear && history -c && history -w" C-m

# select bash pane
tmux select-pane -t $SESSION_NAME:0.1

# Attach to the tmux session
tmux attach-session -t $SESSION_NAME
