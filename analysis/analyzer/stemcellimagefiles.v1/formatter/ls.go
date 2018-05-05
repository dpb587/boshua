package formatter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/dpb587/boshua/analysis/analyzer/stemcellimagefiles.v1/output"
)

type Ls struct{}

func (f Ls) Format(writer io.Writer, reader io.Reader) error {
	var rows [][]string
	var maxlength = []int{0, 0, 0, 0, 0, 0, 0}

	s := bufio.NewScanner(reader)
	for s.Scan() {
		var result output.Result

		err := json.Unmarshal(s.Bytes(), &result)
		if err != nil {
			return err
		}

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
			"-", // TODO try to manually determine link count?
			gname,
			uname,
			strconv.FormatInt(result.Size, 10),
			f.timeText(result),
			result.Path,
		}

		if result.Link != "" {
			row[6] = fmt.Sprintf("%s -> %s", row[6], result.Link)
		}

		rows = append(rows, row)
		if len(row[1]) > maxlength[1] {
			maxlength[1] = len(row[1])
		}
		if len(row[2]) > maxlength[2] {
			maxlength[2] = len(row[2])
		}
		if len(row[3]) > maxlength[3] {
			maxlength[3] = len(row[3])
		}
		if len(row[4]) > maxlength[4] {
			maxlength[4] = len(row[4])
		}
	}

	format := fmt.Sprintf(
		"%%s %%%ds %%%ds %%%ds %%%ds %%s %%s\n",
		maxlength[1],
		maxlength[2],
		maxlength[3],
		maxlength[4],
	)

	for _, row := range rows {
		fmt.Fprintf(writer, format, row[0], row[1], row[2], row[3], row[4], row[5], row[6])
	}

	return nil
}

func (f Ls) timeText(result output.Result) string {
	ts := result.ModTime
	then := time.Now().Add(-1 * time.Second * 86400 * 180)

	if ts.Unix() > then.Unix() {
		return ts.Format("Jan _2 15:04")
	}

	return ts.Format("Jan _2  2006")
}

func (f Ls) modeText(result output.Result) string {
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
