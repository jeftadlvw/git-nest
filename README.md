# git-nest

**git-nest** is a _git_ command line extension for loosely nesting external repositories into your repository. It's an alternative to the default `git-submodule` that does add submodule information to the index, but just works from a configuration file.

> [!IMPORTANT]
> This project is in early development. Source code may be unstable and documentation inaccurate. Proceed and use with caution. Feel free to open an issue!

## Installation

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

`nix-shell`, `devbox shell` and ~~Docker~~ (not yet) are supported out-of-the-box too. More about the dependencies and shell environments can be found in [Development](#development).

The building procedure is defined as a Makefile target.

```shell
make
```

You'll find the compiled binary at `./build/git-nest`. Make sure the binary is in your PATH variable, so _git_ is able to find it.

### Using Docker
> TODO

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
- `debug`: echo project root directory and other Makefile variables

## Roadmap
We are working on it. We have many ideas and much room for improvements. We'll structure and prioritize our internal list before releasing it to the public.

## Contributing
Reviewing and accepting contributions is temporarily suspended until the project's foundation is established. More information will follow soon. Issues and bugs however are welcome!
