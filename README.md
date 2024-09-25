# ESO Tools

(c)2024 Donovan C. Young

This software is published under the terms of the MIT Software License

ESO Tools is a command line tool designed to help manage ESO AddOns

## Installation

### Using Go

The best method to obtain the latest release is directly via Go itself, so
if you have GoLang already installed, you can get the latest release via

```sh
go install github.com/dyoung522/esotools@latest
```

### Direct binary download

Otherwise, you can download the pre-compiled binaries from the Assests section in our [GitHub Releases](https://github.com/dyoung522/esotools/releases)

Feel free to rename them to `esotools` on your particular system for ease of use.

## Configuration

Upon launch, the tool will attempt to auto-discover your ESO_HOME directory. If it cannot, it will prompt you to enter it manually.

This is usually in `My Documents\Elder Scrolls Online` on Windows or `${HOME}/Documents/Elder Scrolls Online` on MacOS / Linux.

Optionally, you can specify the `--esohome` flag, set an `ESO_HOME` environment variable, or create a `.esotools` file in your HOME directory with the following contents:

```yaml
eso_home: "/<your-home-directory>/Documents/Elder Scrolls Online"
```

The advantage to this method is that you only need to do this once, and then can simply run `esotools` without supplying any additional information.

**_PLEASE NOTE: DO NOT include the `live` folder as part of your path._**

## Usage

```sh
Usage:

  esotools [command]


Available Commands:

  backup    Various backup commands
  check     Various check commands
  completion Generate the autocompletion script for the specified shell
  help      Help about any command
  list      Various listing commands


Flags:

      --config string   config file (default is $HOME/.esotools.yaml)
  -H, --esohome live    The full installation path of your ESO game files (where the live folder lives).
  -h, --help            help for esotools
  -N, --no-color        do not output ANSI color codes
  -v, --verbose count   counted verbosity
      --version         version for esotools


Use "esotools [command] --help" for more information about a command.
```

### Commands

#### backup savedvars

```sh
Creates a ZIP backup file of all SavedVariables in the current directory.


Usage:

  esotools backup savedvars [flags]


Flags:

  -h, --help   help for savedvars
```

#### check addons

```sh
Checks AddOns installed in the ESO AddOns directory, and reports any errors


Usage:

  esotools check addons [flags]


Flags:

  -h, --help       help for addons
  -o, --optional   Warn if optional dependencies aren't installed as well
```

#### check savedvars [--backup|--clean|--dryrun]

```sh
Specifically, it reports on extraneous SavedVariable files that do not correspond to any known AddOn.
Optionally, you can auto-remove them with the --clean flag.


Usage:

  esotools check savedvars [flags]


Flags:

      --backup    Performs a backup prior to any destructive actions
      --clean     Removes extranious SavedVariable files
      --dry-run   Shows what changes would be made without actually making them. Use this to double-check before using --clean
  -h, --help      help for savedvars
```

#### list addons

```sh
Lists AddOns installed in the ESO AddOns directory.

By default, this will print out a simple list with only one AddOn per line. However, other formats may be specified via the flags.


Usage:

  esotools list addons [flags]


Flags:

  -h, --help       help for addons
  -j, --json       Print out the list in JSON format
  -m, --markdown   Print out the list in markdown format
  -D, --no-deps    Suppresses printing of AddOns that are dependencies of other AddOns
  -L, --no-libs    Suppresses printing of AddOns that are considered Libraries
  -r, --raw        Print out the list in the RAW ESO AddOn header format (most verbose)
  -s, --simple     Prints the AddOn listing in simple plain text
```
