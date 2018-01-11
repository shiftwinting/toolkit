package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/tideways/toolkit/xhprof"

	"github.com/olekukonko/tablewriter"
)

type Unit struct {
	Name    string
	Divisor float32
}

var (
	ms Unit = Unit{Name: "ms", Divisor: 1000.0}
	kb Unit = Unit{Name: "KB", Divisor: 1024.0}
)

type FieldInfo struct {
	Name   string
	Label  string
	Header string
	Unit   Unit
}

var fieldsMap map[string]FieldInfo = map[string]FieldInfo{
	"wt": FieldInfo{
		Name:   "WallTime",
		Label:  "Inclusive Wall-Time",
		Header: "Wall-Time",
		Unit:   ms,
	},
	"excl_wt": FieldInfo{
		Name:   "ExclusiveWallTime",
		Label:  "Exclusive Wall-Time",
		Header: "Wall-Time",
		Unit:   ms,
	},
	"cpu": FieldInfo{
		Name:   "CpuTime",
		Label:  "Inclusive CPU-Time",
		Header: "CPU-Time",
		Unit:   ms,
	},
	"excl_cpu": FieldInfo{
		Name:   "ExclusiveCpuTime",
		Label:  "Exclusive CPU-Time",
		Header: "CPU-Time",
		Unit:   ms,
	},
	"memory": FieldInfo{
		Name:   "Memory",
		Label:  "Inclusive Memory",
		Header: "Memory",
		Unit:   kb,
	},
	"excl_memory": FieldInfo{
		Name:   "ExclusiveMemory",
		Label:  "Exclusive Memory",
		Header: "Memory",
		Unit:   kb,
	},
	"io": FieldInfo{
		Name:   "IoTime",
		Label:  "Inclusive I/O-Time",
		Header: "I/O-Time",
		Unit:   ms,
	},
	"excl_io": FieldInfo{
		Name:   "ExclusiveIoTime",
		Label:  "Exclusive I/O-Time",
		Header: "I/O-Time",
		Unit:   ms,
	},
}

func renderProfile(profile *xhprof.Profile, field string, fieldInfo FieldInfo, minPercent float32) error {
	profile.SortBy(fieldInfo.Name)
	main, err := profile.GetMain()
	if err != nil {
		return err
	}

	minValue := minPercent * main.GetFloat32Field(fieldInfo.Name)

	var fields []FieldInfo
	if strings.HasPrefix(field, "excl_") {
		fields = []FieldInfo{fieldsMap[strings.TrimPrefix(field, "excl_")], fieldInfo}
	} else {
		fields = []FieldInfo{fieldInfo, fieldsMap["excl_"+field]}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Function", "Count", fieldInfo.Header, fmt.Sprintf("Excl. %s", fieldInfo.Header)})
	for _, call := range profile.Calls {
		if call.GetFloat32Field(fieldInfo.Name) < minValue {
			break
		}

		table.Append(getRow(call, fields))
	}

	fmt.Printf("Showing XHProf data by %s\n", fieldInfo.Label)
	table.Render()

	return nil
}

func getRow(call *xhprof.Call, fields []FieldInfo) []string {
	res := []string{
		fmt.Sprintf("%.90s", call.Name),
		fmt.Sprintf("%d", call.Count),
	}

	for _, field := range fields {
		col := fmt.Sprintf("%2.2f %s", call.GetFloat32Field(field.Name)/field.Unit.Divisor, field.Unit.Name)
		res = append(res, col)
	}

	return res
}
