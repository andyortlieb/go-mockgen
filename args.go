package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/alecthomas/kingpin"
)

const Version = "0.1.0"

var (
	app = kingpin.New("go-mockgen", "go-mockgen generates mock implementations from interface definitions.").Version(Version)

	importPaths    = app.Arg("path", "The import paths used to search for eligible interfaces").Required().Strings()
	pkgName        = app.Flag("package", "The name of the generated package. Is the name of target directory if dirname or filename is supplied by default.").Short('p').String()
	prefix         = app.Flag("prefix", "A prefix used in the name of each mock struct. Should be TitleCase by convention.").String()
	interfaces     = app.Flag("interfaces", "A whitelist of interfaces to generate given the import paths.").Short('i').Strings()
	outputFilename = app.Flag("filename", "The target output file. All mocks are writen to this file.").Short('o').String()
	outputDir      = app.Flag("dirname", "The target output directory. Each mock will be written to a unique file.").Short('d').String()
	force          = app.Flag("force", "Do not abort if a write to disk would overwrite an existing file.").Short('f').Bool()
	listOnly       = app.Flag("list", "Dry run - print the interfaces found in the given import paths.").Bool()
)

var identPattern = regexp.MustCompile("^[A-Za-z]([A-Za-z0-9_]*[A-Za-z])?$")

func parseArgs() (string, string, error) {
	args := os.Args[1:]

	if _, err := app.Parse(args); err != nil {
		return "", "", err
	}

	dirname, filename, err := validateOutputPath(*outputDir, *outputFilename)
	if err != nil {
		return "", "", err
	}

	if *pkgName == "" && !*listOnly {
		if dirname == "" {
			kingpin.Fatalf("could not infer package, try --help")
		}

		*pkgName = path.Base(dirname)
	}

	if !identPattern.Match([]byte(*pkgName)) {
		kingpin.Fatalf("illegal package name supplied, try --help")
	}

	if *prefix != "" && !identPattern.Match([]byte(*prefix)) {
		kingpin.Fatalf("illegal prefix supplied, try --help")
	}

	return dirname, filename, nil
}

func validateOutputPath(dirname, filename string) (string, string, error) {
	if dirname == "" && filename == "" {
		return "", "", nil
	}

	if filename != "" {
		if dirname != "" {
			kingpin.Fatalf("dirname and filename are mutually exclusive, try --help")
		}

		filename, err := filepath.Abs(filename)
		if err != nil {
			return "", "", err
		}

		dirname, filename = path.Dir(filename), path.Base(filename)
	}

	dirname, err := filepath.Abs(dirname)
	if err != nil {
		return "", "", err
	}

	exists, err := pathExists(dirname)
	if err != nil {
		return "", "", err
	}

	if !exists {
		return "", "", fmt.Errorf("directory %s does not exist", dirname)
	}

	return dirname, filename, nil
}
