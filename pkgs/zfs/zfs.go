package zfs

// zfs is a helper function to wrap typical calls to zfs that ignores stdout.
func zfs(arg ...string) error {
	_, err := zfsOutput(arg...)
	return err
}

// zfs is a helper function to wrap typical calls to zfs.
func zfsOutput(arg ...string) ([][]string, error) {
	c := command{Command: "zfs"}
	return c.Run(arg...)
}
