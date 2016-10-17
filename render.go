package main

import (
	"encoding/json"
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
	Ifc  string
	RCs  []*Record
	Cost time.Duration
}

type RecordGroupSortByCost struct{ RecordGroup }
type RecordGroupSortByName struct{ RecordGroup }

func (r RecordGroup) Len() int      { return len(r.RCs) }
func (r RecordGroup) Swap(i, j int) { r.RCs[i], r.RCs[j] = r.RCs[j], r.RCs[i] }
func (r RecordGroup) Sort(s SortBy) {
	switch s {
	case SortByCost:
		sort.Sort(&RecordGroupSortByCost{r})
	case SortByName:
		sort.Sort(&RecordGroupSortByName{r})
	}
}
func (r RecordGroupSortByCost) Less(i, j int) bool {
	oi, oj := r.RCs[i], r.RCs[j]
	return oi.EndAt.Sub(oi.StartAt) < oj.EndAt.Sub(oj.StartAt)
}
func (r RecordGroupSortByName) Less(i, j int) bool {
	oi, oj := r.RCs[i], r.RCs[j]
	return oi.Name < oj.Name
}

type SortRecordGroup []RecordGroup

func (rg SortRecordGroup) Len() int           { return len(rg) }
func (rg SortRecordGroup) Swap(i, j int)      { rg[i], rg[j] = rg[j], rg[i] }
func (rg SortRecordGroup) Less(i, j int) bool { return rg[i].Cost > rg[j].Cost }

func (rg RecordGroup) Since(after time.Time) RecordGroup {
	var ret = RecordGroup{
		Cost: rg.Cost,
		Ifc:  rg.Ifc,
		RCs:  make([]*Record, 0),
	}
	for _, r := range rg.RCs {
		if r.EndAt.After(after) {
			ret.RCs = append(ret.RCs, r)
		}
	}
	return ret
}

func (db Database) Render(w io.Writer, s SortBy, after time.Duration) {
	var data SortRecordGroup
	for _, rg := range db.data {
		rs := rg.Since(db.launchTimestamp.Add(after))
		data = append(data, rs)
	}
	sort.Sort(data)

	json.NewEncoder(w).Encode(data)
}

func (db *Database) RenderGlobalInfo(w io.Writer) {
	var n int
	var cost time.Duration
	for _, rg := range db.data {
		n += len(rg.RCs)
		cost += rg.Cost
	}
	json.NewEncoder(w).Encode(
		struct {
			Cost            time.Duration
			N               int
			LaunchTimestamp time.Time
			BusAddr         string
			Duration        time.Duration
		}{cost, n, db.launchTimestamp, db.busAddr, time.Since(db.launchTimestamp)},
	)
}
