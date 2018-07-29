package analyzer

import (
	"syscall"
	"time"
)

func getFSTimes(stat *syscall.Stat_t) (*time.Time, *time.Time, *time.Time) {
	atime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	mtime := time.Unix(int64(stat.Mtim.Sec), int64(stat.Mtim.Nsec))

	return &atime, &ctime, &mtime
}
