package analyzer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/stemcellimagefiles.v1/output"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/checksum/algorithm"
)

func (a Analyzer) walkFS(results analysis.Writer, baseDir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fullPath := filepath.Join(path, info.Name())
		relPath := fmt.Sprintf("/%s", strings.TrimPrefix(strings.TrimPrefix(fullPath, baseDir), "/"))

		statSys, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			panic("failed to convert stat")
		}

		atime := time.Unix(int64(statSys.Atimespec.Sec), int64(statSys.Atimespec.Nsec))
		ctime := time.Unix(int64(statSys.Ctimespec.Sec), int64(statSys.Ctimespec.Nsec))

		result := output.Result{
			Type: string(info.Mode()),
			Path: relPath,
			Size: statSys.Size,
			Mode: int64(statSys.Mode),
			Uid:  int64(statSys.Uid),
			Gid:  int64(statSys.Gid),
			//			Uname:   header.Uname,
			//			Gname:   header.Gname,
			ModTime:    time.Unix(int64(statSys.Mtimespec.Sec), int64(statSys.Mtimespec.Nsec)),
			AccessTime: &atime,
			ChangeTime: &ctime,
		}

		if info.Mode()&os.ModeSymlink != 0 {
			resolved, err := os.Readlink(fullPath)
			if err != nil {
				return fmt.Errorf("reading link: %v", err)
			}

			result.Link = resolved
		}

		if info.Mode()&os.ModeType == 0 {
			fh, err := os.Open(fullPath)
			if err != nil {
				return fmt.Errorf("open file for checksum: %v", err)
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
				return fmt.Errorf("creating checksum: %v", err)
			}

			result.Checksums = checksums.ImmutableChecksums()
		}

		err = results.Write(result)
		if err != nil {
			return fmt.Errorf("writing result: %v", err)
		}

		return nil
	}
}
