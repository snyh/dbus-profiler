package main

import (
	"fmt"
	"io/ioutil"
	"pkg.deepin.io/lib/dbus"
	"sync"
	"time"
)

type Record struct {
	Name   string
	Ifc    string
	Sender string
	Serial uint32

	StartAt time.Time
	Cost    time.Duration
}

type Database struct {
	sync.RWMutex
	conn       *dbus.Conn
	callerInfo map[uint32]string

	data map[string]RecordGroup

	pending map[uint32]*Record

	busAddr         string
	launchTimestamp time.Time
}

func NewDatabase(bus_addr string, c chan *dbus.Message) (*Database, error) {
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

	db := &Database{
		conn:            conn,
		callerInfo:      make(map[uint32]string),
		data:            make(map[string]RecordGroup),
		pending:         make(map[uint32]*Record),
		busAddr:         bus_addr,
		launchTimestamp: time.Now(),
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

func (rg *RecordGroup) Add(rc *Record) {
	rg.rcs = append(rg.rcs, rc)
	rg.TotalCost += rc.Cost
	rg.TotalCall += 1
}

func (d Database) AddRecord(rc *Record) {
	if rc.Ifc == "org.freedesktop.DBus" {
		return
	}
	rc.Cost = time.Since(rc.StartAt)
	d.Lock()

	ifc := rc.Ifc

	r, ok := d.data[ifc]
	if !ok {
		r.Ifc = ifc
	}
	r.Add(rc)

	d.data[ifc] = r
	d.Unlock()
}

func buildPendingRecord(conn *dbus.Conn, msg *dbus.Message) *Record {
	name, _ := msg.Headers[dbus.FieldMember].Value().(string)
	ifc, _ := msg.Headers[dbus.FieldInterface].Value().(string)
	sender, _ := msg.Headers[dbus.FieldSender].Value().(string)
	switch ifc {
	case "org.freedesktop.DBus.Properties":
		ifc, _ = msg.Body[0].(string)
		switch name {
		case "Get":
			name = "Get(" + msg.Body[1].(string) + ")"
		case "Set":
			name = "Set(" + msg.Body[1].(string) + ")"
		default:
			name = "unknown:(" + name + ")"
		}
	}

	return &Record{
		Name:    name,
		Ifc:     ifc,
		Sender:  sender,
		Serial:  msg.Serial(),
		StartAt: time.Now(),
	}
}

func (db *Database) commitPending(serial uint32) {
	db.RLock()
	rc, ok := db.pending[serial]
	db.RUnlock()

	if ok {
		db.AddRecord(rc)
	}
}

func (db *Database) stashPending(rc *Record) {
	db.Lock()
	db.pending[rc.Serial] = rc
	db.Unlock()
}

func (db *Database) handleInput(msg *dbus.Message) error {
	rc := buildPendingRecord(db.conn, msg)
	switch msg.Type {
	case dbus.TypeMethodCall:
		if msg.Flags&dbus.FlagNoReplyExpected == dbus.FlagNoReplyExpected {
			db.AddRecord(rc)
		} else {
			db.stashPending(rc)
		}
	case dbus.TypeMethodReply:
		rs, _ := msg.Headers[dbus.FieldReplySerial].Value().(uint32)
		db.commitPending(rs)
	case dbus.TypeError:
		rs, _ := msg.Headers[dbus.FieldReplySerial].Value().(uint32)
		db.commitPending(rs)
	case dbus.TypeSignal:
		db.AddRecord(rc)

	default:
		return fmt.Errorf("unknown msg %v\n", msg)
	}
	return nil
}
