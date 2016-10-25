package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Type byte

const (
	TypeMethodCall Type = 1 + iota
	TypeSignal
	TypePropertyGet
	TypePropertySet
	TypePropertyChanged
	typeMax
)

type Record struct {
	Sender string
	OPath  string
	Ifc    string
	Name   string

	Caller string

	Type    Type
	StartAt time.Time
	Cost    time.Duration
}

func (rc Record) Valid() bool {
	if rc.Type == 0 ||
		rc.StartAt.IsZero() ||
		rc.Sender == "" ||
		rc.OPath == "" ||
		rc.Ifc == "" ||
		rc.Name == "" {
		return false
	}
	return true
}

type DatabaseSource interface {
	Source() chan *Record
	QuerySender(string) (*SenderInfo, error)
}

type Database struct {
	sync.RWMutex
	data map[string]RecordGroup

	source          DatabaseSource
	launchTimestamp time.Time
}

func NewDatabase() *Database {
	return &Database{
		data:            make(map[string]RecordGroup),
		launchTimestamp: time.Now(),
	}
}

func (db *Database) QuerySender(s string) *SenderInfo {
	info, err := db.source.QuerySender(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, "QuerySender E:", err)
		return NewErrorSenderInfo(s)
	}
	return info
}
func (db *Database) AddSource(source DatabaseSource) {
	go func() {
		db.source = source
		for r := range source.Source() {
			db.AddRecord(r)
		}
	}()
}

func (d *Database) AddRecord(rc *Record) {
	if !rc.Valid() {
		fmt.Fprintf(os.Stderr, "Invalid Record:%v\n", rc)
		return
	}
	rc.Cost = time.Since(rc.StartAt)

	ifc := rc.Ifc

	d.RLock()
	r, ok := d.data[ifc]
	d.RUnlock()

	if !ok {
		r.Ifc = ifc
	}
	r.Add(rc)

	d.Lock()
	d.data[ifc] = r
	d.Unlock()
}
