package main

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"os"
	"pkg.deepin.io/lib/dbus"
	"time"
)

func (d *Introspector) cacheCaller(sender string) (*SenderInfo, error) {
	if s, ok := d.callerInfo[sender]; ok {
		return s, nil
	} else {
		var pid uint32
		err := d.conn.BusObject().Call("org.freedesktop.DBus.GetConnectionUnixProcessID", 0, sender).Store(&pid)
		if pid == 0 {
			return nil, fmt.Errorf("GetConnectionUnixProcessFail: %v", err)
		}
		info, err := NewSenderInfo(sender, int32(pid))
		if err != nil {
			fmt.Fprintln(os.Stderr, "cacheCaller E:", err, sender)
			return nil, err
		}
		d.callerInfo[sender] = info
		return info, nil
	}
}

func (d *Introspector) Query(sender string) (*SenderInfo, error) {
	return d.cacheCaller(sender)
}

type SenderInfo struct {
	Sender     string
	Pid        int32
	Cmd        string
	CreateTime int64
	End        time.Time
}

func NewErrorSenderInfo(sender string) *SenderInfo {
	return &SenderInfo{
		Sender: sender,
		Pid:    -1,
	}
}

func NewSenderInfo(sender string, pid int32) (*SenderInfo, error) {

	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("NewProcess failed PID:%d E:%v", pid, err)
	}
	cmd, err := p.Cmdline()
	if err != nil {
		return nil, err
	}
	ctime, err := p.CreateTime()
	if err != nil {
		return nil, err
	}

	return &SenderInfo{
		Sender:     sender,
		Pid:        pid,
		Cmd:        cmd,
		CreateTime: ctime,
	}, nil
}

type Introspector struct {
	conn       *dbus.Conn
	callerInfo map[string]*SenderInfo
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
		callerInfo: make(map[string]*SenderInfo),
	}, nil
}
