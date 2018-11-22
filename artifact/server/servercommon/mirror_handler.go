package servercommon

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cheggaaa/pb"
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/metalink/file"
	"github.com/dpb587/boshua/metalink/file/metaurl/boshreleasesource"
	"github.com/dpb587/metalink/file/metaurl"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/pkg/errors"
)

type MirrorHandler struct{}

func (c *MirrorHandler) ServeHTTPArtifact(w http.ResponseWriter, r *http.Request, subject artifact.Artifact) {
	urlLoader := urldefaultloader.New()
	metaurlLoader := metaurl.NewLoaderFactory()
	metaurlLoader.Add(boshreleasesource.Loader{})

	subjectFile := subject.MetalinkFile()

	if subjectFile.Size != 0 {
		w.Header().Set("Content-Size", strconv.FormatUint(subjectFile.Size, 10))
	}

	if subjectFile.Name != "" {
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, subjectFile.Name))
	} else {
		w.Header().Set("Content-Disposition", "attachment")
	}

	for _, hash := range subjectFile.Hashes {
		w.Header().Add("Digest", fmt.Sprintf("%s=%s", hash.Type, base64.StdEncoding.EncodeToString([]byte(hash.Hash))))
	}

	if r.Method == http.MethodHead {
		return
	}

	// TODO refactor; use metaurls; use other urls; UnverifiedTransfer
	progress := pb.New64(0)
	progress.SetWriter(bytes.NewBuffer(nil))

	remote, err := urlLoader.Load(subjectFile.URLs[0])
	if err != nil {
		panic(errors.Wrap(err, "loading remote")) // TODO !panic
	}

	client := file.NewHTTPResponse(w)

	err = client.WriteFrom(remote, progress)
	if err != nil {
		panic(errors.Wrap(err, "transferring file")) // TODO !panic
	}
}
