# asciinemagen

Create an asciinemagen cast file from a flat text file. By default lines
starting with `>` (configurable with `iprompt` flag) are considered command
lines and other lines are considered to be command output (stdout/err, no
distinction for asciinema).

You can control which timing, prompts, and asciinema metadata with flags, see
`asciinemagen --help` for the details.

## Install

At the moment, you'll need to have Go installed, and use `go get` to pull
and build the binary:

```sh
go get github.com/npaton/asciinemagen
```

## Example Usage

```sh
> echo "> date\nTue Mar 13 10:57:27 EST 2018" | asciinemagen --title "Example" > example.cast
```

```sh
> asciinema play example.cast
# Plays cast ...
```

```sh
> asciinema upload example.cast
https://asciinema.org/a/Aag0imG6hcdq6bQ
```
