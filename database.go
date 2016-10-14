package main

import (
	"fmt"
	"io/ioutil"
	"pkg.deepin.io/lib/dbus"
	"time"
)

//type Message struct {
//	Type
//	Flags
//	Headers map[HeaderField]Variant
//	Body []interface{}

//	serial uint32
//}

type Record struct {
	Desc   string
	Begin  time.Time
	Cost   float64
	Sender string
}

type Database struct {
	conn       *dbus.Conn
	callerInfo map[uint32]string
	data       map[string][]*Record
}

func NewDatabase(conn *dbus.Conn, c chan *dbus.Message) (*Database, error) {
	db := &Database{
		conn:       conn,
		callerInfo: make(map[uint32]string),
		data:       make(map[string][]*Record),
	}
	go func() {
		for msg := range c {
			db.handleInput(msg)
		}
		fmt.Println("End!!!")
	}()
	return db, nil
}

func (d Database) QueryCaller(sender uint32) string {
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

func (d Database) AddRecord(rc *Record) {
	fmt.Printf("Add : %v\n", rc)
}

func buildPendingRecord(conn *dbus.Conn, msg *dbus.Message) *Record {
	name, _ := msg.Headers[dbus.FieldMember].Value().(string)
	ifcName, _ := msg.Headers[dbus.FieldInterface].Value().(string)
	sender, _ := msg.Headers[dbus.FieldSender].Value().(string)

	desc := ifcName + "." + name
	switch desc {
	case "org.freedesktop.DBus.Properties.Get":
		desc = "PG: " + msg.Body[0].(string) + "." + msg.Body[1].(string)
	case "org.freedesktop.DBus.Properties.Set":
		desc = "PS: " + msg.Body[0].(string) + "." + msg.Body[1].(string)
	}

	return &Record{
		Desc:   desc,
		Begin:  time.Now(),
		Sender: sender,
	}
}

func (db Database) handleInput(msg *dbus.Message) {
	rc := buildPendingRecord(db.conn, msg)
	db.AddRecord(rc)
}
