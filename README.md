# git-nest

**git-nest** is a _git_ command line extension for nesting external repositories in your project without your parent repository noticing, using native features and configurations files.


> This project is in early development.<br>
> Source code may be unstable and documentation inaccurate.<br>
> Proceed and use with caution!

## Installation

### From source

You should have the following requirements installed into your local environment:
* supported version of the Go toolchain (see [Supported go versions](#supported-go-versions))
* GNU make

`nix-shell`, `devbox shell` and ~~Docker~~ (not yet) are supported out-of-the-box too. More about the dependencies and shell environments can be found in [Development](#development).

The building procedure is written inside the Makefile.

```shell
make
```

You'll find the compiled binary at `./build/git-nest`. Make sure the binary is in your PATH variable, so _git_ is able to find it.

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
> ToDo

### Bare bones

In order to develop *git-nest* without a wrapped shell environment you are required to have some dependencies installed on your system. Some links and documentation are provided.

* supported version of the Go toolchain (see [Supported go versions](#supported-go-versions))
* git
* GNU make

### Supported Go versions
As of current development, `go-1.22` is required to build the source code.

### Makefile targets
- `build` (default): build the project
- `git-test`: test git integration
- `root-dir`: echo project root directory concatenated by make


## Contributing
Reviewing and accepting contributions is temporarily suspended until the project's foundation is established. More information will follow soon.
