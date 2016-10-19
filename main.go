package main

import (
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "dbus-profiler"
	app.Usage = "Profiling dbus message with beautiful, dynamical UI and realtime data"
	app.Version = "0.1"
	app.Author = "snyh"
	app.Email = "snyh@snyh.org"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bus_addr,a",
			Usage: "monitor arbitrary message from the bus ADDRESS [system|user|$dbus_address]",
			Value: "user",
		},
		cli.StringFlag{
			Name:  "server_addr,s",
			Usage: "the address to bind for serving, [127.0.0.1:8080|auto]",
			Value: ":7799",
		},
		cli.BoolFlag{
			Name:  "quiet,q",
			Usage: "disable auto launch web browser",
		},

		cli.StringFlag{
			Name:   "input",
			Usage:  "the cache file path, It will be created by invoke \"socat EXEC:'dbus-monitor --pcap --$bus' PIPE:$VALUE\"",
			Value:  "auto-socat",
			Hidden: true,
		},
		cli.BoolFlag{
			Name:   "debug,d",
			Hidden: true,
		},
	}
	app.Action = func(c *cli.Context) error {
		db := NewDatabase()

		source, err := newPcapSource(c.GlobalString("bus_addr"))
		if err != nil {
			return err
		}
		db.AddSource(source)

		s := NewServer(db, c.GlobalString("server_addr"))

		s.OpenBrowser(c.GlobalBool("quiet"))

		s.Run(c.GlobalBool("debug"))
		if err != nil {
			return err
		}
		return nil
	}
	app.RunAndExitOnError()
}
