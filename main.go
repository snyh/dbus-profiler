package main

// Use tcpdump to create a test file
// tcpdump -w test.pcap
// or use the example above for writing pcap files

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"os"
	"os/exec"
	"pkg.deepin.io/lib/dbus"
	"time"
)

func RawRecord(bus_addr string, input string) (chan *dbus.Message, error) {
	var socat_arg1 string
	switch bus_addr {
	case "user", "session":
		socat_arg1 = "EXEC:'dbus-monitor --session --pcap'"
	case "system":
		socat_arg1 = "EXEC:'dbus-monitor --system --pcap'"
	default:
		socat_arg1 = fmt.Sprintf("EXEC:\"dbus-monitor --addr='%s' --pcap\"", bus_addr)
	}

	switch input {
	case "auto-socat":
		input = fmt.Sprintf("%s/dbus-profiler.%d", os.TempDir(), os.Getpid())
		bin, err := exec.LookPath("socat")
		if err != nil {
			return nil, fmt.Errorf("input=auto-socat need the binary socat be installed on system")
		}
		cmd := exec.Command(bin,
			"-u",
			socat_arg1,
			fmt.Sprintf("PIPE:'%s'", input),
		)
		err = cmd.Start()
		if err != nil {
			return nil, err
		}
		// Wait one second to setup the pipe
		<-time.After(time.Second * 1)
		go cmd.Wait()
	}

	handle, err := pcap.OpenOffline(input)
	if err != nil {
		return nil, err
	}

	//defer	handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	ch := make(chan *dbus.Message, 100)

	go func() {
		for packet := range packetSource.Packets() {
			r := bytes.NewReader(packet.Data())
			msg, err := dbus.DecodeMessage(r)
			if err != nil {
				fmt.Printf("E: %v\n", err)
				continue
			}
			ch <- msg
		}
		close(ch)
	}()
	return ch, nil
}

func main() {
	app := cli.NewApp()
	app.Name = "dbus-profiler"
	app.Usage = "Profiling dbus message with beautiful, dynamical UI and realtime data"
	app.Version = "0.1"
	app.Author = "snyh"
	app.Email = "snyh@snyh.org"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address",
			Usage: "monitor an arbitrary message bus given at ADDRESS [system|user|$dbus_address]",
			Value: "user",
		},
		cli.StringFlag{
			Name:  "input",
			Usage: "the cache file path, It will be created by invoke \"socat EXEC:'dbus-monitor --pcap --$bus' PIPE:$VALUE\"",
			Value: "auto-socat",
		},
		cli.StringFlag{
			Name:  "bind",
			Usage: "the address to bind for serving, [127.0.0.1:8080|auto]",
			Value: ":7799",
		},
		cli.BoolFlag{
			Name:   "debug,d",
			Hidden: true,
		},
		cli.BoolFlag{
			Name:  "quiet,q",
			Usage: "disable auto launch web browser",
		},
	}
	app.Action = func(c *cli.Context) error {
		bus_addr := c.GlobalString("address")

		ch, err := RawRecord(bus_addr, c.GlobalString("input"))
		if err != nil {
			return err
		}

		db, err := NewDatabase(bus_addr, ch)
		if err != nil {
			return err
		}

		s := NewServer(db, c.GlobalString("bind"))
		s.OpenBrowser(c.GlobalBool("quiet"))
		s.Run(c.GlobalBool("debug"))
		if err != nil {
			return err
		}
		return nil
	}
	app.RunAndExitOnError()
}
