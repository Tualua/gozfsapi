package zfs

import (
	"slices"
	"strings"
)

type Dataset struct {
	Name          string
	Origin        string
	Used          uint64
	Avail         uint64
	Mountpoint    string
	Compression   string
	Type          string
	Written       uint64
	Volsize       uint64
	Logicalused   uint64
	Usedbydataset uint64
	Quota         uint64
	Referenced    uint64
}

var (
	// List of ZFS properties to retrieve from zfs list command on a non-Solaris platform.
	dsPropList = []string{"name", "origin", "used", "available", "mountpoint", "compression", "type", "volsize", "quota", "referenced", "written", "logicalused", "usedbydataset"}

	dsPropListOptions = strings.Join(dsPropList, ",")
)

func GetDataset(name string) (*Dataset, error) {
	out, err := zfsOutput("list", "-Hp", "-o", dsPropListOptions, name)
	if err != nil {
		return nil, err
	}
	data := out[0]
	ds := &Dataset{
		Name: data[0],
		Origin: parseString(data[1]),
		Used: parseUint64(data[2]),
		Avail: parseUint64(data[3]),
		Mountpoint: parseString(data[4]),
		Compression: parseString(data[5]),
		Type: parseString(data[6]),
		Volsize: parseUint64(data[7]),
		Quota: parseUint64(data[8]),
		Referenced: parseUint64(data[9]),
		Written: parseUint64(data[10]),
		Logicalused: parseUint64(data[11]),
		Usedbydataset: parseUint64(data[12]),
	}
	return ds, nil
}

func ListDatasets(t string) []string {
	if t == "" {
		t = "all"
	}
	out, err := zfsOutput("list", "-Hp", "-o", "name", "-t", t)
	if err != nil || len(out) == 0 {
		return nil
	}
	return slices.Concat(out...)
}