# Dmenual
A `dmenu` wrapper with very simple, manual configuration.

## Installation
If you have Go installed
```bash
$ go get -u github.com/ravernkoh/dmenual
```
Alternatively, you can download the correct release for your architecture from the [releases](https://github.com/ravernkoh/dmenual/releases) page.

## Configuration
To use `dmenual`, you must first add all the configuration file. The default configuration directory is `~/.config/dmenual` but it can be changed using the `path` flag.

Currently, `dmenual` supports two types of launching: `cli` and `gui`. Executables listed in the `cli` file will be opened in a terminal while those listed in the `gui` file will be opened standalone.

An example configuration file:
```
ranger
htop
python
rtv
```

Comments and whitespace are not supported. Each executable should be seperated by a single newline.

## Usage
Run it like any other executable
```bash
# Opens normal dmenual
$ dmenual

# Change the configuration path
$ dmenual -path ~/custom

# Change the terminal used
$ dmenual -term termite

# Customize the internal dmenu call
$ dmenual -- -fg '#000000'
```
