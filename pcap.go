package main

import (
	"bytes"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"os"
	"os/exec"
	"pkg.deepin.io/lib/dbus"
	"strings"
	"sync"
	"time"
)

type pcapSource struct {
	sync.RWMutex
	pending   map[uint32]*Record
	queue     chan *Record
	cachePath string
	intro     *Introspector
}

func (p *pcapSource) Source() chan *Record {
	return p.queue
}

func newPcapSource(bus_addr string) (*pcapSource, error) {
	intro, err := NewIntrospector(bus_addr)
	if err != nil {
		return nil, err
	}

	bin, err := exec.LookPath("socat")
	if err != nil {
		return nil, fmt.Errorf("You have to install the 'socat' and 'dbus-monitor' for running %q.", os.Args[0])
	}
	source := &pcapSource{
		pending:   make(map[uint32]*Record),
		queue:     make(chan *Record),
		cachePath: fmt.Sprintf("%s/dbus-profiler.%d", os.TempDir(), os.Getpid()),
		intro:     intro,
	}

	var socat_arg1 string
	switch bus_addr {
	case "user", "session":
		socat_arg1 = "EXEC:'dbus-monitor --session --pcap'"
	case "system":
		socat_arg1 = "EXEC:'dbus-monitor --system --pcap'"
	default:
		socat_arg1 = fmt.Sprintf("EXEC:\"dbus-monitor --addr='%s' --pcap\"", bus_addr)
	}
	cmd := exec.Command(bin,
		"-u",
		socat_arg1,
		fmt.Sprintf("PIPE:'%s'", source.cachePath),
	)
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	// Wait one second to setup the pipe
	<-time.After(time.Second * 1)
	go cmd.Wait()
	go source.Run()
	return source, nil
}
func (w *pcapSource) Run() error {
	handle, err := pcap.OpenOffline(w.cachePath)
	if err != nil {
		return err
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		r := bytes.NewReader(packet.Data())
		msg, err := dbus.DecodeMessage(r)
		if err != nil {
			fmt.Printf("E: %v\n", err)
			continue
		}
		w.handleInput(msg)
	}
	close(w.queue)
	return nil
}

func (s *pcapSource) commitPending(serial uint32) {
	s.RLock()
	rc, ok := s.pending[serial]
	s.RUnlock()

	if ok {
		s.queue <- rc
	}
}

func (s *pcapSource) stashPending(serial uint32, rc *Record) {
	s.Lock()
	s.pending[serial] = rc
	s.Unlock()
}

func initRecord(msg *dbus.Message) *Record {
	name, _ := msg.Headers[dbus.FieldMember].Value().(string)
	ifc, _ := msg.Headers[dbus.FieldInterface].Value().(string)
	opath, _ := msg.Headers[dbus.FieldPath].Value().(dbus.ObjectPath)
	sender, _ := msg.Headers[dbus.FieldSender].Value().(string)

	mtype := typeMax

	switch msg.Type {
	case dbus.TypeSignal:
		mtype = TypeSignal
	case dbus.TypeMethodCall:
		mtype = TypeMethodCall
	default:
		panic(fmt.Sprintf("Unknown msg: %v", msg))
	}

	if ifc == "org.freedesktop.DBus.Properties" {
		ifc, _ = msg.Body[0].(string)
		switch name {
		case "Get":
			mtype = TypePropertyGet
			name = "Get(" + msg.Body[1].(string) + ")"

		case "Set":
			mtype = TypePropertySet
			name = "Set(" + msg.Body[1].(string) + ")"

		case "GetAll":
			mtype = TypePropertyGet
			name = "GETALL"
		case "PropertiesChanged":
			mtype = TypePropertyChanged
			cs, ok := msg.Body[1].(map[string]dbus.Variant)
			if !ok {
				fmt.Fprintf(os.Stderr, "invliad ProppertiesChanged msg: %v", msg.Body[1])
			}
			var keys []string
			for k := range cs {
				keys = append(keys, k)
			}
			name = "Change(" + strings.Join(keys, " ") + ")"
		default:
			panic("unknown:(" + name + ")")
		}
	}

	return &Record{
		Sender:  sender,
		OPath:   string(opath),
		Ifc:     ifc,
		Name:    name,
		Type:    mtype,
		StartAt: time.Now(),
	}
}

func (s *pcapSource) handleInput(msg *dbus.Message) error {
	switch msg.Type {
	case dbus.TypeMethodCall:
		rc := initRecord(msg)
		if msg.Flags&dbus.FlagNoReplyExpected == dbus.FlagNoReplyExpected {
			s.send(rc)
		} else {
			s.stashPending(msg.Serial(), rc)
		}
	case dbus.TypeMethodReply, dbus.TypeError:
		replyS, _ := msg.Headers[dbus.FieldReplySerial].Value().(uint32)
		s.commitPending(replyS)
	case dbus.TypeSignal:
		rc := initRecord(msg)
		s.send(rc)

	default:
		return fmt.Errorf("unknown msg %v\n", msg)
	}
	return nil
}

func (s *pcapSource) send(rc *Record) {
	if !rc.Valid() {
		panic("Invalid Record")
	}
	s.queue <- rc
}
