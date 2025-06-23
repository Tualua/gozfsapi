package zfs

import (
	"log"
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

	ds := &Dataset{Name: name}
	log.Println(len(out))
	// for _, line := range out {
	// 	if err := ds.parseLine(line); err != nil {
	// 		return nil, err
	// 	}
	// }

	return ds, nil
}