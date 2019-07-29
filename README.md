# clr - Output Colorizer

`clr` is a command line output colorizer. Output from any command can be piped to `clr` and it'll highlight interesting portions of the output based on the given color rules.

## Installation

1. Clone repo
1. Run `go build`
1. Copy the `clr` binary to `/usr/local/bin`

## Usage

For a detailed guide, hit:
```console
clr -help
```

#### Examples
```console
tail -F log.txt | clr INFO~green ERROR~red
```

## Screenshots

#### Hello, World
<img src="doc/simple.png" width="65%" height="65%"/>

#### Regex and more
<img src="doc/complex.png" width="65%" height="65%"/>