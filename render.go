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

type RecordSummary struct {
	Ifc       string
	TotalCost time.Duration
	TotalCall int

	CostDetail []time.Duration
	CallDetail []int
}

type CallUsageSummary struct {
	Total int
	Cost  []time.Duration
}

type RecordDetail struct {
	Ifc       string
	TotalCost time.Duration
	TotalCall int

	Method   map[string]CallUsageSummary
	Signal   map[string]CallUsageSummary
	Property map[string]CallUsageSummary
}

func calcUsage(v []*Record) CallUsageSummary {
	if len(v) == 0 {
		panic("Zero v in calcUsage")
	}

	cost := make([]time.Duration, 0)
	if v[0].Type != TypeSignal {
		for _, rc := range v {
			cost = append(cost, rc.Cost)
		}
	}
	return CallUsageSummary{
		Cost:  cost,
		Total: len(v),
	}
}

func (rg RecordGroup) Detail() RecordDetail {
	mcache := make(map[string][]*Record)
	scache := make(map[string][]*Record)
	pcache := make(map[string][]*Record)

	for _, rc := range rg.rcs {
		switch rc.Type {
		case TypeMethodCall:
			mcache[rc.Name] = append(mcache[rc.Name], rc)
		case TypeSignal:
			scache[rc.Name] = append(scache[rc.Name], rc)
		case TypePropertyGet, TypePropertySet:
			pcache[rc.Name] = append(pcache[rc.Name], rc)
		}
	}
	//TODO: Move the logic of above in RecordGroup and Database.AddRecord

	ret := RecordDetail{
		Ifc:       rg.Ifc,
		TotalCost: rg.TotalCost,
		TotalCall: rg.TotalCall,
		Method:    make(map[string]CallUsageSummary),
		Signal:    make(map[string]CallUsageSummary),
		Property:  make(map[string]CallUsageSummary),
	}
	for n, v := range mcache {
		ret.Method[n] = calcUsage(v)
	}
	for n, v := range scache {
		ret.Signal[n] = calcUsage(v)
	}
	for n, v := range pcache {
		ret.Property[n] = calcUsage(v)
	}
	return ret
}

type RecordGroup struct {
	Ifc       string
	TotalCost time.Duration
	TotalCall int
	rcs       []*Record
}

type SortRecordGroup []RecordSummary

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

func (rg RecordGroup) Summary(after time.Time, unit time.Duration) RecordSummary {
	var ret = RecordSummary{
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

func (db *Database) Render(w io.Writer, top int, last time.Duration) {
	ts := db.launchTimestamp

	since := time.Since(db.launchTimestamp)

	if since > last {
		ts = db.launchTimestamp.Add(since - last)
	}

	var ret []RecordSummary
	db.RLock()
	for _, rg := range db.data {
		ret = append(ret, rg.Summary(ts, time.Second))
	}
	db.RUnlock()

	sort.Sort(SortRecordGroup(ret))

	if top < len(ret) {
		ret = ret[0:top]
	}
	json.NewEncoder(w).Encode(ret)
}

func (db *Database) RenderInterfaceDetail(name string, w io.Writer) error {
	db.RLock()
	v, ok := db.data[name]
	db.RUnlock()
	if !ok {
		return fmt.Errorf("There hasn't any record for %s", name)
	}
	return json.NewEncoder(w).Encode(v.Detail())
}

func (db *Database) RenderGlobalInfo(w io.Writer) {
	var n int
	var cost time.Duration
	db.RLock()
	for _, rg := range db.data {
		n += len(rg.rcs)
		cost += rg.TotalCost
	}
	db.RUnlock()
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
