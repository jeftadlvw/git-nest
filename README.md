# git-nest

**git-nest** is a _git_ command line extension for loosely nesting external repositories into your repository. It's an alternative to the default `git-submodule` that doesn't add submodule information to the index, but just works from a configuration file.

> [!IMPORTANT]
> This project is in early development. Source code may be unstable and documentation inaccurate. Proceed and use with caution. Feel free to open an issue!

> [!NOTE]
> I need to put my time and effort into uni projects. This project is not abandoned, but it'll rest for a little bit. Critical bugs for the latest release will be considered, features and other improvements need to wait. Thank you for your patience!


## Installation

UNIX-based operating systems:
```shell
curl -fsSL https://raw.githubusercontent.com/jeftadlvw/git-nest/main/install.sh | bash
```

Windows users:
```shell
powershell -c "irm https://raw.githubusercontent.com/jeftadlvw/git-nest/main/install.ps1 | iex"
```

These scripts also provide a command that adds the installation directory to the user's PATH. Restart your shell afterward.

### Binaries
We provide pre-built binary files for major platforms and architectures at every release. You can find them [here](https://github.com/jeftadlvw/git-nest/releases).

Download the binary for your operating system and processor, store it somewhere safe and rename it to `git-nest`. Add the path you downloaded `git-nest` into to your system `PATH` variables in order for `git` to automatically register it as a subcommand.

Verify you have a working installation by running
```shell
git nest info
```

We plan on providing an installation script.

### From source

You should have the following requirements installed into your local environment:
* supported version of the Go toolchain (see [Supported go versions](#supported-go-versions))
* GNU make

`nix-shell`, `devbox shell` are supported development environments too. More about the dependencies and shell environments can be found in [Development](#development). There also is a Docker image that spins up a [testing environment](#temporary-testing-environment) with automatic code building and other tools.

The building procedure is defined as a Makefile target.

```shell
make
```

You'll find the compiled binary at `_build/git-nest`. Make sure the binary is in your PATH variable, so _git_ is able to find it.

### Using Docker
There also is a Dockerfile that creates an image with which you can create binaries for a specific target.
First, build the image.

```shell
docker build -t git-nest/build _docker/build
```

Then run it and mount the source code in read-only mode to `/git-nest/src` and output directory to `/git-nest/build`:
```shell
docker run \
    -it \
    --rm \
    -v ./:/git-nest/src:ro \
    -v ./_build:/git-nest/build \
    git-nest/build
```

This will build binaries for supported default targets (Windows, Linux, Darwin (MacOS) | amd64, arm64) and generate a `checksums.txt` containing pre-calculated sha256 hashes for the artifacts.

#### Custom build targets
You can specify your own targets with the `TARGET_OVERRIDE` environment variable:
```shell
docker run \
    -it \
    --rm \
    -e TARGET_OVERRIDE="openbsd/amd64" \
    -v ./:/git-nest/src:ro \
    -v ./_build:/git-nest/build \
    git-nest/build
```

This will now build for the openbsd operating system on amd64-bit processors. You can also define multiple targets at once by just separating them using a whitespace, e.g.:
```shell
-e TARGET_OVERRIDE="js/wasm openbsd/amd64 android/arm64"
```

Supported build targets are listed [here](https://go.dev/doc/install/source#environment) (section `$GOOS and $GOARCH`).

#### Injecting version and commit hash
Official release binaries have an injected binary version and commit hash. You can do that too by using the `GIT_NEST_BUILD_VERSION` and `GIT_NEST_BUILD_COMMIT_SHA` environment variables:
```shell
docker run \
    -it \
    --rm \
    -e GIT_NEST_BUILD_VERSION="your-version-tag" \
    -e GIT_NEST_BUILD_COMMIT_SHA="your-commit-hash" \
    -v ./:/git-nest/src:ro \
    -v ./_build:/git-nest/build \
    git-nest/build
```

#### Skipping tests
The script runs all tests before actually creating the binary artifacts. Although not recommended, you can skip the tests using the `SKIP_TESTS` environment variable. It's just required that the variable has a value.
```shell
-e SKIP_TESTS="true" # value can be anything
```

## Getting started
There currently is not much documentation outside from the cli help and this README:
```
$ git nest
Usage:
  git-nest [flags]
  git-nest [command]

Available Commands:
  add         Add and clone a remote submodule into this project
  help        Help about any command
  info        Print various debug information
  list        List nested modules
  pull        Pull new updates in all nested modules
  remove      Remove a submodule from this project
  sync        Update and apply state changes
  verify      Verify configuration and nested modules
  version     Print git-nest version

Flags:
  -h, --help      help for git-nest
  -v, --version   version for git-nest

Use "git-nest [command] --help" for more information about a command.
```
Some quick notes:
- the most relevant commands are `add`, `remove` and `sync`.
- running these commands will create a `nestmodules.toml` file, which hold all the important information about your nested modules. Commit and share this file. See issue #4 ([click](https://github.com/jeftadlvw/git-nest/issues/4#issue-2229919243)) for information on the general structure.
- synchronization between the configuration and existing modules is currently as follows: If the directory does not exist, then the module is created as defined in the configuration file. If the module already exists, every change within the module (like branches, commits, ...) are synchronized into the configuration file. This behaviour is suspect to change within the upcoming releases.

## Development

### Shell environments
We support some wrappers out-of-the-box. In case you are interested in those projects, we also link some documentation.

#### nix
```shell
nix-shell
```

- https://nixos.org/
- https://wiki.nixos.org/wiki/Development_environment_with_nix-shell

#### devbox
```shell
devbox shell
```

- https://www.jetpack.io/devbox
- https://www.jetpack.io/devbox/docs/

### Docker
> ToDo<br>
> Using a [testing environment](#temporary-testing-environment) might be an alternative until we put the whole development environment into a docker image.

### Bare bones
In order to develop *git-nest* without a wrapped shell environment you are required to have some dependencies installed on your system. Some links and documentation are provided.

* supported version of the Go toolchain (see [Supported go versions](#supported-go-versions))
* git
* GNU make

### Temporary testing environment
We provide a Docker image that sets up an isolated testing environment to securely test modified code in a shell.

To build the image, enter
```shell
docker build -t git-nest/test-env _docker/test_env
```

Run it with
```shell
docker run -it --name git-nest_testenv -v ./:/app/src:ro git-nest/test-env
```

In case you want to continue testing
```shell
docker start -ai git-nest_testenv
```

There also is a makefile target that combines these two commands. It will reuse a pre-existing container.
```shell
make test-env
```

This opens a tmux session with two terminals. On the left side you see the output of the file watcher that automatically builds your source code on file changes. On the right side is your regular bash terminal with which you can interact and test _git-nest_ with. There is a command called `prune`, that completely wipes the test-env directory. We also provide some default repositories you can clone. List them with `list_repos`.

### Supported Go versions
As of current development, `go-1.22` is required to build the source code.

### Makefile targets
- `build` (default): build the project
- `git-test`: test git integration
- `debug`: echo project root directory and other Makefile variables
- `test-env`: build and spin up a temporary testing environment

## Roadmap
We are working on it. We have many ideas and much room for improvements. We'll structure and prioritize our internal list before releasing it to the public.

You can have a look at the issues where we curate short-term planned features and improvements.

## Contributing
Reviewing and accepting contributions is temporarily suspended until the project's foundation is established. More information will follow soon. Issues and bugs however are welcome!
