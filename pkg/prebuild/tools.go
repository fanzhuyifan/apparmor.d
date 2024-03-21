// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"strings"

	"github.com/arduino/go-paths-helper"
)

func copyTo(src *paths.Path, dst *paths.Path) error {
	files, err := src.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories(), paths.FilterOutNames("README.md"))
	if err != nil {
		return err
	}
	for _, file := range files {
		destination, err := file.RelFrom(src)
		if err != nil {
			return err
		}
		destination = dst.JoinPath(destination)
		if err := destination.Parent().MkdirAll(); err != nil {
			return err
		}
		if err := file.CopyTo(destination); err != nil {
			return err
		}
	}
	return nil
}

// Overwrite upstream profile: rename our profile & hide upstream
func debianOverwrite(files []string) {
	const ext = ".apparmor.d"
	file, err := paths.New("debian/apparmor.d.hide").Append()
	if err != nil {
		panic(err)
	}
	for _, name := range files {
		origin := RootApparmord.Join(name)
		dest := RootApparmord.Join(name + ext)
		if err := origin.Rename(dest); err != nil {
			panic(err)
		}
		if _, err := file.WriteString("/etc/apparmor.d/" + name + "\n"); err != nil {
			panic(err)
		}
	}
}

// Clean the debian/apparmor.d.hide file
func debianOverwriteClean() {
	const debianHide = `# This file is generated by "make", all edit will be lost.

/etc/apparmor.d/usr.bin.firefox
/etc/apparmor.d/usr.sbin.cups-browsed
/etc/apparmor.d/usr.sbin.cupsd
/etc/apparmor.d/usr.sbin.rsyslogd
`
	path := paths.New("debian/apparmor.d.hide")
	if err := path.WriteFile([]byte(debianHide)); err != nil {
		panic(err)
	}
}

// Get the list of upstream profiles to overwrite from dist/overwrite
func getOverwriteProfiles() []string {
	res := []string{}
	lines, err := DistDir.Join("overwrite").ReadFileAsLines()
	if err != nil {
		panic(err)
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		res = append(res, line)
	}
	return res
}

