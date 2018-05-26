package analyzer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/output"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/boshua/util/checksum/algorithm"
	"github.com/pkg/errors"
)

func (a Analyzer) loadFileNameMap(path string) (map[int]string, error) {
	results := map[int]string{}

	pathBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading file")
	}

	for _, line := range strings.SplitN(string(pathBytes), "\n", -1) {
		cols := strings.SplitN(line, ":", 4)
		if len(cols) != 4 {
			continue
		}

		num, err := strconv.Atoi(cols[2])
		if err != nil {
			return nil, errors.Wrap(err, "converting id")
		}

		results[num] = cols[0]
	}

	return results, nil
}

func (a Analyzer) walkFS(results analysis.Writer, baseDir string, userMap map[int]string, groupMap map[int]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath := fmt.Sprintf("/%s", strings.TrimPrefix(strings.TrimPrefix(path, baseDir), "/"))

		statSys, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			panic("failed to convert stat")
		}

		atime, ctime, mtime := getFSTimes(statSys)

		result := output.Result{
			Type:       string(info.Mode().String()[0]),
			Path:       relPath,
			Size:       statSys.Size,
			Mode:       int64(statSys.Mode),
			Uid:        int64(statSys.Uid),
			Gid:        int64(statSys.Gid),
			Uname:      userMap[int(statSys.Uid)],
			Gname:      groupMap[int(statSys.Gid)],
			ModTime:    *mtime,
			AccessTime: atime,
			ChangeTime: ctime,
		}

		if info.Mode()&os.ModeSymlink != 0 {
			resolved, err := os.Readlink(path)
			if err != nil {
				return errors.Wrap(err, "reading link")
			}

			result.Link = resolved
		}

		if info.Mode()&os.ModeType == 0 {
			fh, err := os.Open(path)
			if err != nil {
				return errors.Wrap(err, "open file for checksum")
			}

			defer fh.Close()

			checksums := checksum.WritableChecksums{
				checksum.New(algorithm.MustLookupName(algorithm.MD5)),
				checksum.New(algorithm.MustLookupName(algorithm.SHA1)),
				checksum.New(algorithm.MustLookupName(algorithm.SHA256)),
				checksum.New(algorithm.MustLookupName(algorithm.SHA512)),
			}

			_, err = io.Copy(checksums, fh)
			if err != nil {
				return errors.Wrap(err, "creating checksum")
			}

			result.Checksums = checksums.ImmutableChecksums()
		}

		err = results.Write(result)
		if err != nil {
			return errors.Wrap(err, "writing result")
		}

		return nil
	}
}
