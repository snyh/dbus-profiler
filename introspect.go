package main

import (
	"fmt"
	"io/ioutil"
	"pkg.deepin.io/lib/dbus"
)

func (d Introspector) QueryCaller(sender uint32) string {
	if s, ok := d.callerInfo[sender]; ok {
		return s
	} else {

		var pid uint32
		d.conn.BusObject().Call("org.freedesktop.DBus.GetConnectionUnixProcessID", 0, sender).Store(&pid)
		if pid == 0 {
			return "unknown"
		}
		info, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
		if err != nil {
			return fmt.Sprintf("PID:%d", pid)
		}
		return fmt.Sprintf("%s(%d)", info, pid)
	}
}

type Introspector struct {
	conn       *dbus.Conn
	callerInfo map[uint32]string
}

func NewIntrospector(bus_addr string) (*Introspector, error) {
	var conn *dbus.Conn
	var err error
	switch bus_addr {
	case "user", "session":
		conn, err = dbus.SessionBus()
	case "system":
		conn, err = dbus.SystemBus()
	default:
		conn, err = dbus.Dial(bus_addr)
	}
	if err != nil {
		return nil, err
	}
	return &Introspector{
		conn:       conn,
		callerInfo: make(map[uint32]string),
	}, nil
}
