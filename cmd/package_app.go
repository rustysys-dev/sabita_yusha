package main

import (
	"fmt"

	"github.com/google/rpmpack"
	"github.com/rustysys-dev/rpmbuild"
)

var app = rpmbuild.Builder{
	BinDir:  "bin",
	DistDir: "dist",

	RPMMetaData: rpmpack.RPMMetaData{
		Name:        "sabita_yusha",
		Summary:     "Configurable `Programmable Button` execution daemon for VIA/QMK keyboards.",
		Description: "runs programmable buttons for VIA/QMK keyboards via user config file. It is primarily built using the yushakobo Quick Paint keyboard, and while i imagine it should work with other keyboards compliant to QMK on Linux I cannot guarantee it.",
		Version:     "1.0.0",
		Release:     "1",
		Arch:        "x86_64",
		Packager:    "Scott Mattan <scott.mattan@rustysys.dev>",
		Licence:     "MIT",
		Compressor:  "zstd",
		Provides: []*rpmpack.Relation{{
			Name:    "sabita_yusha",
			Version: "1.0.0",
		}},
	},

	Files: []rpmbuild.PackageFile{
		{
			Source:      "bin/sabita_yusha",
			Destination: "/usr/bin/sabita_yusha",
		},
		{
			Source:      "scripts/systemd/sabita_yusha.service",
			Destination: "/usr/lib/systemd/user/sabita_yusha.service",
		},
	},
}

func main() {
	if err := app.Build(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	if err := app.Package(); err != nil {
		panic(err)
	}
}
