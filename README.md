# paintstripper - Removes bash/shell color codes from command/log output

A useful command for when you have command or log output with bash/shell color codes like `^ESC[31m` and the shell/interface that is displaying the output does not interpret these correctly (or you simply don't want the output with color).

## Usage

You can use paintstripper either to process existing output files, or better yet, pipe a command you need to clean up directly to paintstripper.

e.g. `./my_colorful_command.sh | paintstripper -color`, this will take the output of `my_colorful_command.sh` and strip the color information away and print to stdout.

For more information, run `paintstripper -help` for usage and some more examples.

## Output control

A useful table to understand what information is displayed (and where defined, logged) based on the arguments used.

| -write-stripped | -write-unstripped | Console output | Other output                                  |
|-----------------|-------------------|----------------|-----------------------------------------------|
|                 |                   | Stripped       | Unstripped output discarded                   |
|       YES       |                   | Unstripped     | Stripped output written to file               |
|                 |        YES        | Stripped       | Unstripped output written to file             |
|       YES       |        YES        | None           | Stripped & snstripped output written to files |

*NOTE: any use of `-quiet` will surpress output to stdout. Use `-force` to overwrite any existing files when logging.*

## Installation

Either [download the latest release](https://github.com/sparkcodeuk/paintstripper/releases) for your OS and stick the `paintstripper` binary (making sure it's executable) in your `$PATH`, or build the release yourself.

To build the release clone the repo and run the `.../paintstripper/util/build-release.sh` shell script. This will run some tests and then build for all supported OS's their respective binaries. You can then grab the one you need. (You can also run `go build`, but you'll forgo the tests that `build-release.sh` runs).

## Future improvements

* Functional tests
* Stripping all control characters (not just shell color escape codes)
* Carriage return processing (e.g., when used to output a progress bar etc.)
