package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"
)

type SortBy int

const (
	SortByCost = iota
	SortByName
)

type RecordGroup struct {
	Ifc       string
	TotalCost time.Duration
	TotalCall int
	rcs       []*Record
}

type SortRecordGroup []RecordDetail

func (rg SortRecordGroup) Len() int           { return len(rg) }
func (rg SortRecordGroup) Swap(i, j int)      { rg[i], rg[j] = rg[j], rg[i] }
func (rg SortRecordGroup) Less(i, j int) bool { return rg[i].TotalCost > rg[j].TotalCost }

func (rg RecordGroup) reduce(after time.Time, before time.Time) RecordGroup {
	var ret = RecordGroup{
		TotalCost: rg.TotalCost,
		TotalCall: rg.TotalCall,
		Ifc:       rg.Ifc,
		rcs:       make([]*Record, 0),
	}
	for _, r := range rg.rcs {
		if r.StartAt.After(after) && r.StartAt.Before(before) {
			ret.rcs = append(ret.rcs, r)
		}
	}
	return ret
}

func (rg RecordGroup) Detail(after time.Time, unit time.Duration) RecordDetail {
	var ret = RecordDetail{
		Ifc:        rg.Ifc,
		TotalCost:  rg.TotalCost,
		TotalCall:  rg.TotalCall,
		CallDetail: make([]int, 0),
		CostDetail: make([]time.Duration, 0),
	}

	for {
		before := after.Add(unit)
		t := rg.reduce(after, before)
		ret.CallDetail = append(ret.CallDetail, t.CurrentCall())
		ret.CostDetail = append(ret.CostDetail, t.CurrentCost())
		after = before
		if before.After(time.Now()) {
			break
		}
	}
	return ret
}

func (rg RecordGroup) CurrentCall() int {
	return len(rg.rcs)
}
func (rg RecordGroup) CurrentCost() time.Duration {
	var c time.Duration
	for _, rc := range rg.rcs {
		c += rc.Cost
	}
	return c
}

type RecordDetail struct {
	Ifc        string
	TotalCost  time.Duration
	TotalCall  int
	CostDetail []time.Duration
	CallDetail []int
}

func (db *Database) Render(w io.Writer, top int, last time.Duration) {
	ts := db.launchTimestamp

	since := time.Since(db.launchTimestamp)

	if since > last {
		ts = db.launchTimestamp.Add(since - last)
	}

	var ret []RecordDetail
	for _, rg := range db.data {
		ret = append(ret, rg.Detail(ts, time.Second))
	}

	sort.Sort(SortRecordGroup(ret))

	if top < len(ret) {
		ret = ret[0:top]
	}
	json.NewEncoder(w).Encode(ret)
}

func (db *Database) RenderInterface(name string, w io.Writer) error {
	v, ok := db.data[name]
	if !ok {
		return fmt.Errorf("There hasn't any record for %s", name)
	}

	return json.NewEncoder(w).Encode(v)
}

func (db *Database) RenderGlobalInfo(w io.Writer) {
	var n int
	var cost time.Duration
	for _, rg := range db.data {
		n += len(rg.rcs)
		cost += rg.TotalCost
	}
	json.NewEncoder(w).Encode(
		struct {
			Cost            time.Duration
			N               int
			LaunchTimestamp time.Time
			BusAddr         string
			Duration        time.Duration
		}{cost, n, db.launchTimestamp, "user", time.Since(db.launchTimestamp)},
	)
}

func (rg *RecordGroup) Add(rc *Record) {
	rg.rcs = append(rg.rcs, rc)
	rg.TotalCost += rc.Cost
	rg.TotalCall += 1
}
