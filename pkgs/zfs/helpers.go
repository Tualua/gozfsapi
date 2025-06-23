package zfs

import "strconv"

func parseUint64(s string) uint64 {
	var value uint64 = 0
	if s != "-" {
		value, _ = strconv.ParseUint(s, 10, 64)
	}
	return value
}

func parseString(s string) string {
	if s == "-" {
		return ""
	}
	return s
}
