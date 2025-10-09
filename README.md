# km-probe
`km-probe` is an utility to help you find common mistakes within your Karaoke Mugen repositories.

## When km-probe could be useful
- As a repository maintainer, you can run `km-probe` find potential mistakes in all non-commited new karaokes and karaokes modifications (for example, karaokes with resolution not set to 0×0).
- You can run `km-probe` on all karaoke of a repository to find some karaoke having a specific mistake that you want to work on (for example karaoke without automation script).

Not yet implemented, but could be done eventually:
- Goal for v1.0.0: Integration directly into Karaoke Mugen application as an optional plugin.
- Goal for v2.0.0: As a contributor, you could get an early feedback on karaoke you sent to inbox before a maintainer get the time to look at it (and you may want to update your inbox because you find a mistake).

## What km-probe is not
- `km-probe` is not a replacement for maintainers. Some checks can only be manual and will never be automated.
- `km-probe` is not intended to force contributors or maintainers to create karaoke according to "the only correct way to do karaokes".
It is only a tool, and it can be wrong. Maintainers still have final word on what is or is not integrated into their karaoke repository,
and should rely on their own judgment to do this.
- `km-probe` is not bug-free. If you find a bug, please [report it](https://github.com/karaoke-tools/km-probe/issues) if it is not already.
- `km-probe` is not mandatory. You can run Karaoke Mugen application without `km-probe`.
However, I hope it will help you making better karaokes and following your repository rules.

## Supported platforms
I intend to support the same set of platform Karaoke Mugen does by the release of v1.0.0.
Currently, has only be tested on Debian 12, not extensively.

If you have issues running it on your platform, please [report it](https://github.com/karaoke-tools/km-probe/issues).
Please note that I don't have any MacOS or Windows license, so I cannot test on those platforms myself,
but I will treat your platform specific issues with all the required care.

If you tested successfuly on a different platform, it would be helping to inform me (via the [Github discussions feature](https://github.com/karaoke-tools/km-probe/discussions) for example).

## Semantic Versioning
This software is using [Semantic Versioning](https://semver.org/). The public API will be fully defined by the release of v1.0.0.

Until v1.0.0, we are still in the initial development phase, so backward compatibility is not yet ensured.

## How to use
To get a list of available commands and options, run `km-probe help`.

### List available probes
To list available probes, run `km-probe info`.

### Run on karaokes
If you want to run probes on all karaokes, you can use the command `km-probe karaokes --all`

### Run only on a subset of karaokes
You can select precisely which karaokes to analyse by providing a KID (Karaoke UUID) with the `--kid` flag.

For example:

```bash
$ km-probe karaokes --kid=c8e61289-31fd-430a-8a56-f0ed95f84d50
```

This flag can be provided multiple times to analyse several karaokes.


### Run only on new/modified karaokes
If you want to run only on new/modified karaoke (in comparison to the last local commit on the repository), you can run the command `km-probe git`.

### Global options
#### Repositories
By default, all enabled repositories are searched. If you want to search only from a subset of enabled repositories, you can use the flag `--repo`.
When this flag is used with a repository name, all other repositories will be considered disabled.

For example:

```bash
$ km-probe --repo=mugen.re karaokes --all
```

This flag can be provided multiple times to enable several repositories.

#### Output format
By default, the output format is detected automatically depending on the context.
If the output is in a terminal, it will be displayed in human readable format.
If the output is redirected to a file or another tool, the output is using `json` format.
This allow you to easily re-use the output with other tools like [`jq`](https://github.com/jqlang/jq).

You can force the output to one of those format using the global `--output-format` option with `txt` or `json` as value.

#### Colors
By default, human output will be displayed using ANSI escape codes for color.
Some old terminal may not support ANSI escape codes. You can disable color using the `--color=never` option, or with the environment variable [`NO_COLOR`](https://no-color.org/).

#### Hyperlinks
By default, an hyperlink is created on the songname of analysed karaokes.
It targets a modification page from your Karaoke Mugen application for this karaoke.
Unfortunately, it is currently not possible to open links directly into the application instead of in your browser (but it may be in the future).

Not all terminals support hyperlinks. Some requires a special configuration to enable them but as long as ANSI escaping codes are supported
by your terminal, you should not notice them even if there is no specific support for hyperlinks.

If hyperlinks annoy you, they can be disabled using the `--hyperlink=never` option.

## How to install
### From release artifacts
1. Download the binary that matches your system for the latest release.
2. In all this documentation we will assume you have renamed it to `km-probe`
3. Make it executable (on Linux `chmod +x km-probe`)
4. Move it to a directory present in your `PATH` environment variable (so you can use the software without indicating the full path).
For example, on Linux `mkdir -p ~/bin && mv km-probe ~/bin` should work on most distributions.

### Other methods
For all the methods below, you must have Go toolchain installed on your computer.
You can use the package provided by your distribution (recommended) or [download the latest version](https://go.dev/dl/).
You must use a version of Go above 1.21
(and it will automatically download the toolchain version used by this project).

#### Install using golang builtin package manager
Install the latest version:
```bash
$ go install github.com/karaoke-tools/km-probe@latest
```

or install a specific version (replace `vX.Y.Z` by the version number):

```bash
$ go install github.com/karaoke-tools/km-probe@vX.Y.Z
```


#### From source
##### Build
0. Dependencies: `git`, `make` (they are most likely already installed on your system)
1. Clone the repository (`git clone https://github.com/karaoke-tools/km-probe`)
2. Change directory to go inside the cloned repository (`cd km-probe`)
2. Run: `make`, this creates the binary as `km-probe`

##### Install
Run: `sudo make install`

With this last method, bash auto-completions will be installed automatically.

##### Uninstall
Run: `sudo make uninstall`

##### Update
To update the software, run `git pull`, and do the Build and the Install steps again.

## License
`km-probe` is a free and open source software.
See the [`LICENSE`](https://github.com/karaoke-tools/km-probe/blob/master/LICENSE) file for more information.
