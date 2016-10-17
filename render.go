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
		sort.Sort(RecordGroupSortByCost{r})
	case SortByName:
		sort.Sort(RecordGroupSortByName{r})
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

func (db *Database) Render(w io.Writer, s SortBy) {
	var keys []string
	for k := range db.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var data []RecordGroup
	for _, key := range keys {
		rs := db.data[key]
		rs.Sort(s)
		data = append(data, rs)
	}

	json.NewEncoder(w).Encode(data)
}

type GlobalInfo struct {
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
