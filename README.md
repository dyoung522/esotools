# ESO Tools

(c)2024 Donovan C. Young

This software is published under the terms of the MIT Software License

ESO Tools is a command line tool designed to help manage ESO AddOns

## Installation

Installation instructions are coming Soon...

## Configuration

Before you can use the tool, it needs to know the location where your Local ESO game configuration files are installed.
This is usually in `%AppData%\Elder Scrolls Online` on Windows or `${HOME}/Documents/Elder Scrolls Online` on MacOS.

### Configuration File

The best way to do this is to create an `.esotools.yaml` file in your HOME directory with the following contents:

```yaml
eso_home: "/<your-home-directory>/Documents/Elder Scrolls Online"
```

The advantage to creating an `.esotools.yaml` file is that you only need to do this once, and then can simply run `esotools` without supplying any additional information.

### Other Configuration Options

Optionally, you can also use the command line option `--esohome` (or `-H` for short), or set an `ESO_HOME` environment variable. However, please note that using either of these options will override anything in your `.esotools.yaml` file (if it exists).

#### command-line option

```sh
esotools -H "${HOME}/Documents/Elder Scrolls Online"
```

#### environment variable

```sh
ESO_HOME="${HOME}/Documents/Game Files/Elder Scrolls Online" esotools
```
