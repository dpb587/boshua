package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/dpb587/boshua/metalink/metalinkutil"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/pkg/errors"
	// releaseversionv2 "github.com/dpb587/boshua/releaseversion/api/v2/server"
	// stemcellversionv2 "github.com/dpb587/boshua/stemcellversion/api/v2/server"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type BOSHioV2 struct {
	releaseIndex  releaseversiondatastore.Index
	stemcellIndex stemcellversiondatastore.Index
	logger        logrus.FieldLogger
}

func NewBOSHioV2(
	logger logrus.FieldLogger,
	releaseIndex releaseversiondatastore.Index,
	stemcellIndex stemcellversiondatastore.Index,
) *BOSHioV2 {
	return &BOSHioV2{
		releaseIndex:  releaseIndex,
		stemcellIndex: stemcellIndex,
		logger:        logger.WithField("build.package", reflect.TypeOf(BOSHioV2{}).PkgPath()),
	}
}

type boshioV1Releases []boshioV1Release

type boshioV1Release struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url"`
	SHA1    string `json:"sha1"`
}

type boshioV1Stemcells []boshioV1Stemcell

type boshioV1Stemcell struct {
	Name    string                `json:"name"`
	Version string                `json:"version"`
	Regular boshioV1StemcellAsset `json:"regular"`
	Light   boshioV1StemcellAsset `json:"light,omitempty"`
}

type boshioV1StemcellAsset struct {
	URL  string `json:"url"`
	MD5  string `json:"md5"`
	SHA1 string `json:"sha1"`
}

func (h *BOSHioV2) Mount(m *mux.Router) {
	m.HandleFunc("/d/{repo:.+}", func(w http.ResponseWriter, r *http.Request) {
		repoPath := mux.Vars(r)["repo"]
		if repoPath == "" {
			// TODO !panic
			panic(errors.New("empty repository"))
		}

		version := r.FormValue("v")
		if version == "" {
			// TODO !panic
			panic(errors.New("empty version"))
		}

		filterParams := releaseversiondatastore.FilterParams{
			VersionExpected: true,
			Version:         version,
			LabelsExpected:  true,
			Labels: []string{
				fmt.Sprintf("repo/%s", repoPath),
			},
		}

		artifacts, err := h.releaseIndex.GetArtifacts(filterParams, releaseversiondatastore.SingleArtifactLimitParams)
		if err != nil {
			// TODO !panic
			panic(errors.Wrap(err, "getting artifact"))
		}

		// TODO bosh.io adds s3 headers to force downloaded file name; currently uses raw bucket object name
		http.Redirect(w, r, artifacts[0].SourceTarball.URLs[0].URL, http.StatusFound)
	})

	m.HandleFunc("/api/v1/releases/{repo:.+}", func(w http.ResponseWriter, r *http.Request) {
		repoPath := mux.Vars(r)["repo"]
		if repoPath == "" {
			// TODO !panic
			panic(errors.New("empty repository"))
		}

		filterParams := releaseversiondatastore.FilterParams{
			LabelsExpected: true,
			Labels: []string{
				fmt.Sprintf("repo/%s", repoPath),
			},
		}

		// TODO respect all=(1|true|t)
		limitParams := releaseversiondatastore.LimitParams{
			LimitExpected: true,
			Limit:         40,
		}

		artifacts, err := h.releaseIndex.GetArtifacts(filterParams, limitParams)
		if err != nil {
			// TODO !panic
			panic(errors.Wrap(err, "getting artifacts"))
		}

		var result boshioV1Releases

		for _, artifact := range artifacts {
			sha1, err := metalinkutil.HashesToChecksums(artifact.SourceTarball.Hashes).GetByAlgorithm("sha1")
			if err != nil {
				// TODO log warning?
				continue
			}

			// TODO verify absolute url
			downloadURL := url.URL{Scheme: r.URL.Scheme, Host: r.URL.Host, Path: fmt.Sprintf("/d/%s", repoPath), RawQuery: fmt.Sprintf("v=%s", artifact.Version)}

			result = append(
				result,
				boshioV1Release{
					Name:    repoPath, // unfortunate
					Version: artifact.Version,
					URL:     downloadURL.String(),
					SHA1:    strings.TrimPrefix(sha1.String(), "sha1:"),
				},
			)
		}

		// responseBytes, err := json.MarshalIndent(result, "", "  ")
		responseBytes, err := json.Marshal(result)
		if err != nil {
			panic(err) // TODO !panic
		}

		w.Write(responseBytes)
		w.Write([]byte("\n"))
	})

	m.HandleFunc("/api/v1/stemcells/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := mux.Vars(r)["slug"]
		if slug == "" {
			// TODO !panic
			panic(errors.New("empty slug"))
		}

		filterParams := stemcellversiondatastore.FilterParamsFromSlug(slug)

		// TODO merge light + heavy results
		filterParams.FlavorExpected = true
		filterParams.Flavor = "heavy"

		// TODO respect all=(1|true|t)
		limitParams := stemcellversiondatastore.LimitParams{
			LimitExpected: true,
			Limit:         40,
		}

		artifacts, err := h.stemcellIndex.GetArtifacts(filterParams, limitParams)
		if err != nil {
			// TODO !panic
			panic(errors.Wrap(err, "getting artifacts"))
		}

		var result boshioV1Stemcells

		for _, artifact := range artifacts {
			sha1, err := metalinkutil.HashesToChecksums(artifact.MetalinkFile().Hashes).GetByAlgorithm("sha1")
			if err != nil {
				// TODO log warning?
				continue
			}

			result = append(
				result,
				boshioV1Stemcell{
					Name:    artifact.FullName(),
					Version: artifact.Version,
					Regular: boshioV1StemcellAsset{
						// TODO validate url
						URL:  artifact.MetalinkFile().URLs[0].URL,
						SHA1: strings.TrimPrefix(sha1.String(), "sha1:"),
					},
				},
			)
		}

		// responseBytes, err := json.MarshalIndent(result, "", "  ")
		responseBytes, err := json.Marshal(result)
		if err != nil {
			panic(err) // TODO !panic
		}

		w.Write(responseBytes)
		w.Write([]byte("\n"))
	})
}
