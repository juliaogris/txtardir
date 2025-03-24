# txtardir

`txtardir` creates a `txtar` archive from a directory, respecting `.gitignore`
rules.

## Installation

	go install github.com/juliaogris/txtardir@latest

## Usage

	Usage: txtardir [<in> [<out>]] [flags]

	txtar creates a `txtar` archive from a directory, respecting `.gitignore` rules or from a config file.

	Arguments:
	  [<in>]     input directory
	  [<out>]    output file, defaults to stdout

	Flags:
	  -h, --help             Show context-sensitive help.
	  -c, --config=STRING    File containing new line separated list of paths to archived, ignores in directory
