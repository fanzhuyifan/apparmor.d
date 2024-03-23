// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"reflect"
	"testing"

	"github.com/arduino/go-paths-helper"
)

func TestDirective_Usage(t *testing.T) {
	tests := []struct {
		name        string
		d           Directive
		wantMessage string
		wantUsage   string
	}{
		{
			name:        "empty",
			d:           Directives["stack"],
			wantMessage: "Stack directive applied",
			wantUsage:   `#aa:stack profiles_name...`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Usage(); got != tt.wantUsage {
				t.Errorf("Directive.Usage() = %v, want %v", got, tt.wantUsage)
			}
			if got := tt.d.Message(); got != tt.wantMessage {
				t.Errorf("Directive.Usage() = %v, want %v", got, tt.wantMessage)
			}
		})
	}
}

func TestNewOption(t *testing.T) {
	tests := []struct {
		name  string
		file  *paths.Path
		match []string
		want  *Option
	}{
		{
			name: "dbus",
			file: nil,
			match: []string{
				"  #aa:dbus own bus=system name=org.gnome.DisplayManager",
				"dbus",
				"own bus=system name=org.gnome.DisplayManager",
			},
			want: &Option{
				Name: "dbus",
				ArgMap: map[string]string{
					"bus":  "system",
					"name": "org.gnome.DisplayManager",
					"own":  "",
				},
				ArgList: []string{"own", "bus=system", "name=org.gnome.DisplayManager"},
				File:    nil,
				Raw:     "  #aa:dbus own bus=system name=org.gnome.DisplayManager",
			},
		},
		{
			name: "only",
			file: nil,
			match: []string{
				"    #aa:only opensuse",
				"only",
				"opensuse",
			},
			want: &Option{
				Name:    "only",
				ArgMap:  map[string]string{"opensuse": ""},
				ArgList: []string{"opensuse"},
				File:    nil,
				Raw:     "    #aa:only opensuse",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOption(tt.file, tt.match); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		file    *paths.Path
		profile string
		want    string
	}{
		{
			name:    "none",
			file:    nil,
			profile: `  `,
			want:    `  `,
		},
		{
			name:    "present",
			file:    nil,
			profile: `  #aa:dbus own bus=system name=org.freedesktop.systemd1`,
			want:    dbusOwnSystemd1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Run(tt.file, tt.profile); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
