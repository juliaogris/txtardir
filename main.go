package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	ignore "github.com/sabhiram/go-gitignore"
	"golang.org/x/tools/txtar"
)

const description = "txtar creates a `txtar` archive from a directory, respecting `.gitignore` rules."

type app struct {
	In  string `arg:"" help:"input directory" default:"."`
	Out string `arg:"" help:"output file, defaults to stdout" default:""`
}

type shouldSkipFunc func(relPath string) bool

func main() {
	opts := []kong.Option{
		kong.Description(description),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true}),
	}
	kctx := kong.Parse(&app{}, opts...)
	kctx.FatalIfErrorf(kctx.Run())
}

func (a *app) Run() error {
	shouldSkipFn, err := gitignoredSkipper(filepath.Join(a.In, ".gitignore"))
	if err != nil {
		return err
	}

	archive, err := createArchive(a.In, shouldSkipFn)
	if err != nil {
		return err
	}

	return writeArchive(archive, a.Out)
}

func gitignoredSkipper(gitignorePath string) (shouldSkipFunc, error) {
	var matcher ignore.IgnoreParser
	if _, err := os.Stat(gitignorePath); err == nil {
		matcher, err = ignore.CompileIgnoreFile(gitignorePath)
		if err != nil {
			return nil, fmt.Errorf("error parsing .gitignore: %w", err)
		}
	}
	return func(relPath string) bool { return matcher != nil && matcher.MatchesPath(relPath) }, nil
}

func createArchive(inDir string, shouldSkipFn shouldSkipFunc) (*txtar.Archive, error) {
	archive := &txtar.Archive{}

	err := filepath.Walk(inDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil // skip directories
		}

		relPath, err := filepath.Rel(inDir, path)
		if err != nil {
			return err
		}

		if shouldSkipFn(relPath) {
			return nil // skip files according to shouldSkipFn
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		file := txtar.File{
			Name: relPath,
			Data: content,
		}
		archive.Files = append(archive.Files, file)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}
	return archive, nil
}

func writeArchive(archive *txtar.Archive, out string) error {
	w := io.Writer(os.Stdout)
	if out != "" {
		var err error
		w, err = os.Create(out)
		if err != nil {
			return fmt.Errorf("error creating output file: %w", err)
		}
		defer w.(io.Closer).Close()
	}

	outputData := txtar.Format(archive)
	if _, err := w.Write(outputData); err != nil {
		return fmt.Errorf("error writing txtar data: %w", err)
	}

	return nil
}
