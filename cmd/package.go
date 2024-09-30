package main

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/goreleaser/nfpm/v2"
)

var (
	dirs  = []string{"pkg", "scripts"}
	files = []string{"go.mod", "go.sum", "main.go"}
)

type TarWriter struct {
	*tar.Writer
}

func (tw *TarWriter) AddFile(filename string) error {
	// Open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory strucuture would
	// not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = filename

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// if err := os.RemoveAll("dist"); err != nil {
	// 	panic(err)
	// }
	// if err := os.MkdirAll("dist", os.ModePerm); err != nil {
	// 	panic(err)
	// }
	// out, err := os.Create("dist/sabita_yusha.tar")
	// if err != nil {
	// 	panic(err)
	// }
	// defer out.Close()
	//
	// tw := &TarWriter{tar.NewWriter(out)}
	// defer tw.Close()
	//
	// for _, dir := range dirs {
	// 	if err := tw.AddFS(os.DirFS(dir)); err != nil {
	// 		panic(err)
	// 	}
	// }
	//
	// for _, file := range files {
	// 	tw.AddFile(file)
	// }

	createPackage("", "", "")
}

func createPackage(cfg, target, pkger string) error {
	targetIsADirectory := false
	stat, err := os.Stat(target)
	if err == nil && stat.IsDir() {
		targetIsADirectory = true
	}

	if pkger == "" {
		ext := filepath.Ext(target)
		if targetIsADirectory || ext == "" {
			return errors.New("a packager must be specified if target is a directory or blank")
		}

		pkger = ext[1:]
		fmt.Println("guessing packager from target file extension...")
	}

	config, err := nfpm.ParseFile(cfg)
	if err != nil {
		return err
	}

	info, err := config.Get(pkger)
	if err != nil {
		return err
	}

	info = nfpm.WithDefaults(info)

	fmt.Printf("using %s packager...\n", pkger)
	pkg, err := nfpm.Get(pkger)
	if err != nil {
		return err
	}

	if target == "" {
		// if no target was specified create a package in
		// current directory with a conventional file name
		target = pkg.ConventionalFileName(info)
	} else if targetIsADirectory {
		// if a directory was specified as target, create
		// a package with conventional file name there
		target = path.Join(target, pkg.ConventionalFileName(info))
	}

	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	info.Target = target

	if err := pkg.Package(info, f); err != nil {
		os.Remove(target)
		return err
	}

	fmt.Printf("created package: %s\n", target)
	return f.Close()
}
