package metalink

import (
	"bytes"
	"io"

	"github.com/cheggaaa/pb"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/metalink/file"
	"github.com/dpb587/boshua/metalink/file/metaurl/boshreleasesource"
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file/metaurl"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/pkg/errors"
)

func StreamFile(meta4File metalink.File, w io.WriteCloser) error {
	bl := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(bl)

	urlLoader := urldefaultloader.New(fs)
	metaurlLoader := metaurl.NewLoaderFactory()
	metaurlLoader.Add(boshreleasesource.Loader{})

	// TODO refactor; use metaurls; use other urls; UnverifiedTransfer
	progress := pb.New64(0)
	progress.SetWriter(bytes.NewBuffer(nil))

	remote, err := urlLoader.Load(meta4File.URLs[0])
	if err != nil {
		return errors.Wrap(err, "loading remote")
	}

	client := file.NewWriter(w)

	err = client.WriteFrom(remote, progress)
	if err != nil {
		return errors.Wrap(err, "transferring file")
	}

	return nil
}
