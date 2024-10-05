package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/google/rpmpack"
)

var app = builder{
	binDir:  "bin",
	distDir: "dist",

	RPMMetaData: rpmpack.RPMMetaData{
		Name:        "sabita_yusha",
		Summary:     "Configurable `Programmable Button` execution daemon for VIA/QMK keyboards.",
		Description: "runs programmable buttons for VIA/QMK keyboards via user config file. It is primarily built using the yushakobo Quick Paint keyboard, and while i imagine it should work with other keyboards compliant to QMK on Linux I cannot guarantee it.",
		Version:     "1.0.0",
		Arch:        "x86_64",
		Packager:    "Scott Mattan <scott.mattan@rustysys.dev>",
		Licence:     "MIT",
		Compressor:  "zstd",
		Provides: []*rpmpack.Relation{{
			Name:    "sabita_yusha",
			Version: "1.0.0",
		}},
	},

	files: []TargetFile{
		{
			Path:          "build/sabita_yusha",
			InstallTarget: "/usr/bin/sabita_yusha",
		},
		{
			Path:          "scripts/systemd/sabita_yusha.service",
			InstallTarget: "/usr/lib/systemd/user/sabita_yusha.service",
		},
	},
}

type TargetFile struct {
	Path          string
	InstallTarget string
}

// TODO: finish
func (f TargetFile) ToRPMFile() (*rpmpack.RPMFile, error) {
	return nil, nil
}

type builder struct {
	rpmpack.RPMMetaData
	binDir  string
	distDir string
	files   []TargetFile
}

func (b builder) genBinName() string {
	return b.binDir + "/" + b.Name
}

func (b builder) genRPMName() string {
	return b.distDir + "/" + b.Name + ".rpm"
}

func (b builder) Build() error {
	if err := os.RemoveAll("build"); err != nil {
		return err
	}
	if err := os.MkdirAll("build", os.ModePerm); err != nil {
		return err
	}

	stdout, err := exec.Command("go", "build", "-o", b.genBinName(), ".").Output()
	if err != nil {
		fmt.Println(stdout)
		return err
	}

	return nil
}

func (b builder) Package() error {
	if err := os.RemoveAll(b.distDir); err != nil {
		return err
	}

	if err := os.MkdirAll(b.distDir, os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(b.genRPMName())
	if err != nil {
		return err
	}
	defer out.Close()

	r, err := rpmpack.NewRPM(app.RPMMetaData)
	if err != nil {
		return err
	}

	for _, file := range b.files {
		f, err := file.ToRPMFile()
		if err != nil {
			return err
		}

		if f != nil {
			r.AddFile(*f)
		}
	}

	// TODO: need to verify before write?
	if err := r.Write(out); err != nil {
		return err
	}

	return nil
}

// TODO: create a base builder package
func main() {
	if err := app.Build(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	if err := app.Package(); err != nil {
		panic(err)
	}
}
