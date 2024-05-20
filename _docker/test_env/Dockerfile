FROM golang:1.22.3-bookworm

# define directories
ENV GIT_NEST_BASE="/app"
ENV GIT_NEST_SRC="$GIT_NEST_BASE/src"
ENV GIT_NEST_BIN="$GIT_NEST_BASE/bin"
ENV GIT_NEST_TEST_ENV="$GIT_NEST_BASE/test-env"

# create directories
RUN mkdir -p $GIT_NEST_BASE
RUN mkdir -p $GIT_NEST_SRC
RUN mkdir -p $GIT_NEST_BIN
RUN mkdir -p $GIT_NEST_TEST_ENV

# install required packages
RUN apt update
RUN apt install git inotify-tools tmux tree nano -y

# choose base directory and add required files
WORKDIR $GIT_NEST_BASE

# add application entry point
COPY entrypoint.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

# set workdir to test environment directory
WORKDIR $GIT_NEST_TEST_ENV

# create non-root user and give him access to the test env directory
RUN useradd -r -m -s /bin/bash testenv
RUN chown -R testenv:testenv $GIT_NEST_TEST_ENV
RUN chmod -R o+rw $GIT_NEST_BIN
RUN chmod -R o+rw /tmp
USER testenv:testenv

# set container entrypoint
ENTRYPOINT "$GIT_NEST_BASE/entrypoint.sh" $GIT_NEST_SRC $GIT_NEST_BIN $GIT_NEST_TEST_ENV
