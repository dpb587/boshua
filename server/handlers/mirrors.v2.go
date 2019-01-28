package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/cheggaaa/pb"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/metalink/file"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/transfer"
	"github.com/pkg/errors"
	// releaseversionv2 "github.com/dpb587/boshua/releaseversion/api/v2/server"
	// stemcellversionv2 "github.com/dpb587/boshua/stemcellversion/api/v2/server"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MirrorsV2 struct {
	releaseIndex                releaseversiondatastore.Index
	releaseAnalysisIndexGetter  analysisdatastore.NamedGetter
	releaseCompilationIndex     compilationdatastore.Index
	stemcellIndex               stemcellversiondatastore.Index
	stemcellAnalysisIndexGetter analysisdatastore.NamedGetter
	logger                      logrus.FieldLogger
	downloader                  transfer.Transfer
}

func NewMirrorsV2(
	logger logrus.FieldLogger,
	releaseIndex releaseversiondatastore.Index,
	releaseAnalysisIndexGetter analysisdatastore.NamedGetter,
	releaseCompilationIndex compilationdatastore.Index,
	stemcellIndex stemcellversiondatastore.Index,
	stemcellAnalysisIndexGetter analysisdatastore.NamedGetter,
	downloader transfer.Transfer,
) *MirrorsV2 {
	return &MirrorsV2{
		releaseIndex:                releaseIndex,
		releaseCompilationIndex:     releaseCompilationIndex,
		releaseAnalysisIndexGetter:  releaseAnalysisIndexGetter,
		stemcellIndex:               stemcellIndex,
		stemcellAnalysisIndexGetter: stemcellAnalysisIndexGetter,
		downloader:                  downloader,
		logger:                      logger.WithField("build.package", reflect.TypeOf(MirrorsV2{}).PkgPath()),
	}
}

func (h *MirrorsV2) Mount(m *mux.Router) {
	m.HandleFunc("/api/v2/download/release", func(w http.ResponseWriter, r *http.Request) {
		filterParams, err := releaseversiondatastore.FilterParamsFromURLValues(r.URL.Query())
		if err != nil {
			// TODO !panic
			panic(err)
		}

		artifacts, err := h.releaseIndex.GetArtifacts(filterParams, releaseversiondatastore.SingleArtifactLimitParams)
		if err != nil {
			// TODO !panic
			panic(errors.Wrap(err, "getting release"))
		}

		h.stream(w, r, artifacts[0].SourceTarball)
	}).Methods(http.MethodGet)

	m.HandleFunc("/api/v2/download/release-compilation", func(w http.ResponseWriter, r *http.Request) {
		releaseFilterParams, err := releaseversiondatastore.FilterParamsFromURLValues(r.URL.Query())
		if err != nil {
			// TODO !panic
			panic(err)
		}

		osFilterParams, err := osversiondatastore.FilterParamsFromURLValues(r.URL.Query())
		if err != nil {
			// TODO !panic
			panic(err)
		}

		filterParams := compilationdatastore.FilterParams{
			Release: releaseFilterParams,
			OS:      osFilterParams,
		}

		artifacts, err := h.releaseCompilationIndex.GetCompilationArtifacts(
			filterParams,
			// TODO limitparams!
			// releaseversiondatastore.SingleArtifactLimitParams,
		)
		if err != nil {
			// TODO !panic
			panic(errors.Wrap(err, "getting release compilation"))
		}

		h.stream(w, r, artifacts[0].Tarball)
	}).Methods(http.MethodGet)

	m.HandleFunc("/api/v2/download/stemcell", func(w http.ResponseWriter, r *http.Request) {
		filterParams, err := stemcellversiondatastore.FilterParamsFromURLValues(r.URL.Query())
		if err != nil {
			// TODO !panic
			panic(err)
		}

		artifacts, err := h.stemcellIndex.GetArtifacts(filterParams, stemcellversiondatastore.SingleArtifactLimitParams)
		if err != nil {
			// TODO !panic
			panic(errors.Wrap(err, "getting stemcell"))
		}

		h.stream(w, r, artifacts[0].MetalinkFile())
	}).Methods(http.MethodGet)
}

func (h *MirrorsV2) stream(w http.ResponseWriter, r *http.Request, meta4File metalink.File) error {
	progress := pb.New64(int64(meta4File.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)
	progress.SetWriter(bytes.NewBuffer(nil))

	if meta4File.Size > 0 {
		w.Header().Set("content-length", strconv.FormatUint(meta4File.Size, 10))
	}

	if meta4File.Name != "" {
		w.Header().Set("content-disposition", fmt.Sprintf(`attachment; filename="%s"`, url.QueryEscape(meta4File.Name)))
	} else {
		w.Header().Set("content-disposition", fmt.Sprintf("attachment"))
	}

	for _, hash := range meta4File.Hashes {
		w.Header().Add("digest", fmt.Sprintf("%s=%s", hash.Type, base64.StdEncoding.EncodeToString([]byte(hash.Hash))))
	}

	err := h.downloader.TransferFile(
		meta4File,
		file.NewHTTPResponse(w),
		progress,
		nil,
	)
	if err != nil {
		// TODO shouldn't ignore the error; it currently errors because it tries
		// to reopen the stream to verify checksum/signature which fails. should
		// consider how to do an unverified transfer/stream-only
	}

	return nil
}
