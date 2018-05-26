package formatter

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/output"
)

type Ls struct {
	writer    io.Writer
	rows      [][]string
	maxlength []int
}

func NewLs(writer io.Writer) *Ls {
	return &Ls{
		writer:    writer,
		maxlength: []int{0, 0, 0, 0, 0, 0, 0},
	}
}

func (f *Ls) Add(result output.Result) {
	gname := result.Gname
	if gname == "" {
		gname = strconv.FormatInt(result.Gid, 10)
	}

	uname := result.Uname
	if uname == "" {
		uname = strconv.FormatInt(result.Uid, 10)
	}

	row := []string{
		f.modeText(result),
		"0", // TODO try to manually determine link count?
		gname,
		uname,
		strconv.FormatInt(result.Size, 10),
		f.timeText(result),
		result.Path,
	}

	if result.Link != "" {
		row[6] = fmt.Sprintf("%s -> %s", row[6], result.Link)
	}

	f.rows = append(f.rows, row)
	if len(row[1]) > f.maxlength[1] {
		f.maxlength[1] = len(row[1])
	}
	if len(row[2]) > f.maxlength[2] {
		f.maxlength[2] = len(row[2])
	}
	if len(row[3]) > f.maxlength[3] {
		f.maxlength[3] = len(row[3])
	}
	if len(row[4]) > f.maxlength[4] {
		f.maxlength[4] = len(row[4])
	}
}

func (f *Ls) Flush() {
	format := fmt.Sprintf(
		"%%s %%%ds %%%ds %%%ds %%%ds %%s %%s\n",
		f.maxlength[1],
		f.maxlength[2],
		f.maxlength[3],
		f.maxlength[4],
	)

	for _, row := range f.rows {
		fmt.Fprintf(f.writer, format, row[0], row[1], row[2], row[3], row[4], row[5], row[6])
	}
}

func (f *Ls) timeText(result output.Result) string {
	ts := result.ModTime
	then := time.Now().Add(-1 * time.Second * 86400 * 180)

	if ts.Unix() > then.Unix() {
		return ts.Format("Jan _2 15:04")
	}

	return ts.Format("Jan _2  2006")
}

func (f *Ls) modeText(result output.Result) string {
	var buf [10]byte // Mode is uint32.

	if len(result.Type) > 0 {
		buf[0] = result.Type[0]

		// https://golang.org/src/os/types.go?s=1131:1151#L42
		if buf[0] == 'L' {
			buf[0] = 'l'
		} else if buf[0] == 'S' {
			buf[0] = 's'
		} else if buf[0] == 'D' {
			buf[0] = 'b'
		}
	} else {
		buf[0] = '-'
	}

	w := 1
	const rwx = "rwxrwxrwx"
	for i, c := range rwx {
		if result.Mode&(1<<uint(9-1-i)) != 0 {
			buf[w] = byte(c)
		} else {
			buf[w] = '-'
		}
		w++
	}
	return string(buf[:w])
}
