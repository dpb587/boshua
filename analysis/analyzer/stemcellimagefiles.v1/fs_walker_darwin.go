package analyzer

import (
	"syscall"
	"time"
)

func getFSTimes(stat *syscall.Stat_t) (*time.Time, *time.Time, *time.Time) {
	atime := time.Unix(int64(stat.Atimespec.Sec), int64(stat.Atimespec.Nsec))
	ctime := time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec))
	mtime := time.Unix(int64(stat.Mtimespec.Sec), int64(stat.Mtimespec.Nsec))

	return &atime, &ctime, &mtime
}
