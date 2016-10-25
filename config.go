package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type Config struct {
	c string
	l string
}

func NewConfig() Config {
	cP := os.ExpandEnv("$HOME/.config/systemd/user/dbus-profiler.service")
	lP := os.ExpandEnv("$HOME/.config/systemd/user/default.target.wants/dbus-profiler.service")
	return Config{cP, lP}
}

func (c Config) Install() error {
	var self_exec string
	if path.IsAbs(os.Args[0]) {
		self_exec = os.Args[0]
	} else {
		p := path.Base(os.Args[0])
		var err error
		self_exec, err = exec.LookPath(p)
		if err != nil {
			return err
		}
	}

	template := `
[Unit]
Description=dbus-profiler
After=dbus.service

[Service]
ExecStart=` + self_exec + ` -q

[Install]
WantedBy=default.target
`
	return ioutil.WriteFile(c.c, ([]byte)(template), 0755)
}

func (c Config) Enable(v bool) error {
	if v {
		if err := c.Install(); err != nil {
			return err
		}
		return exec.Command("systemctl", "--user", "enable", "dbus-profiler.service").Run()
	} else {
		return exec.Command("systemctl", "--user", "disable", "dbus-profiler.service").Run()
	}
}

func (c Config) CheckEnable() bool {
	l, err := os.Readlink(c.l)
	if err != nil {
		return false
	}

	if l == c.c {
		return true
	}
	return false
}
