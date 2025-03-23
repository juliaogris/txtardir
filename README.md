# txtardir

`txtardir` creates a `txtar` archive from a directory, respecting `.gitignore`
rules.

## Installation

	go install github.com/juliaogris/txtardir@latest

## Usage

	Usage: txtardir [<in> [<out>]]

	txtar is a tool to create a txtar archive from a directory. It ignores files according to a top-level .gitignore file if present

	Arguments:
	  [<in>]     input directory
	  [<out>]    output file, defaults to stdout

	Flags:
	  -h, --help    Show context-sensitive help.
