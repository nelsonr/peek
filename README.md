# peek

Simple tool to preview multiple files on the terminal, built with
[Go](https://go.dev).

![](preview.gif)

## How to install

1. Make sure [Go](https://go.dev) is installed on your system
2. Run the command to install it directly from the repo.

```
go install github.com/nelsonr/peek@latest
```

Alternativately, clone the repo and built directly from the source with:

```
go build
```

## How to use

Run the command `peek` anywhere to preview the files of the current directory.

You also combine with other commands, via the pipe operator.

Example:

```
find . *.log | peek
```

This would create a preview of all log files found in the current directory (and
descending levels).

## Development

Peek was mainly built with the [tview](https://pkg.go.dev/github.com/rivo/tview)
Go package.
