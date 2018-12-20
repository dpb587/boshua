package manifest

import (
	"fmt"
	"strings"

	"github.com/cppforlife/go-patch/patch"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
)

type ReleasePatch struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`

	Source   ReleasePatchRef
	Compiled ReleasePatchRef
	Stemcell Stemcell

	pointer patch.Pointer
}

func (r ReleasePatch) Slug() string {
	return fmt.Sprintf("%s/%s", r.Name, r.Version)
}

func (r ReleasePatch) IsCompiled() bool {
	return r.Compiled.URL != ""
}

func (r ReleasePatch) FilterParams() compilationdatastore.FilterParams {
	f := compilationdatastore.FilterParams{
		Release: releaseversiondatastore.FilterParams{
			NameExpected:    true,
			Name:            r.Name,
			VersionExpected: true,
			Version:         r.Version,
		},
		OS: osversiondatastore.FilterParams{
			NameExpected:    true,
			Name:            r.Stemcell.OS,
			VersionExpected: true,
			Version:         r.Stemcell.Version,
		},
	}

	if r.Source.Sha1 != "" {
		f.Release.ChecksumExpected = true
		f.Release.Checksum = fmt.Sprintf("sha1:%s", r.Source.Sha1)
	}

	if r.Source.URL != "" {
		f.Release.URIExpected = true
		f.Release.URI = r.Source.URL
	}

	return f
}

type ReleasePatchRef struct {
	Sha1 string `yaml:"sha1"`
	URL  string `yaml:"url"`
}

type Stemcell struct {
	OS      string `yaml:"os"`
	Version string `yaml:"version"`
}

func (s Stemcell) Slug() string {
	return fmt.Sprintf("%s/%s", s.OS, s.Version)
}

func (r ReleasePatch) Op() patch.Op {
	if strings.HasSuffix(r.pointer.String(), "/-") {
		value := map[string]interface{}{
			"name":    r.Name,
			"version": r.Version,
		}

		if r.Compiled.URL != "" {
			value["url"] = r.Compiled.URL
			value["sha1"] = strings.TrimPrefix(r.Compiled.Sha1, "sha1:")
			value["stemcell"] = map[string]interface{}{
				"os":      r.Stemcell.OS,
				"version": r.Stemcell.Version,
			}
		} else {
			value["url"] = r.Source.URL
			value["sha1"] = strings.TrimPrefix(r.Source.Sha1, "sha1:")
		}

		return patch.ReplaceOp{
			Path:  r.pointer,
			Value: value,
		}
	}

	ops := patch.Ops{}

	if r.Compiled.URL == "" {
		ops = append(
			ops,
			patch.ReplaceOp{
				Path:  patch.MustNewPointerFromString(fmt.Sprintf("%s/%s", r.pointer.String(), "url?")),
				Value: r.Source.URL,
			},
			patch.ReplaceOp{
				Path:  patch.MustNewPointerFromString(fmt.Sprintf("%s/%s", r.pointer.String(), "sha1?")),
				Value: strings.TrimPrefix(r.Source.Sha1, "sha1:"),
			},
			patch.RemoveOp{
				Path: patch.MustNewPointerFromString(fmt.Sprintf("%s/%s", r.pointer.String(), "stemcell?")),
			},
		)
	} else {
		ops = append(
			ops,
			patch.ReplaceOp{
				Path:  patch.MustNewPointerFromString(fmt.Sprintf("%s/%s", r.pointer.String(), "url?")),
				Value: r.Compiled.URL,
			},
			patch.ReplaceOp{
				Path:  patch.MustNewPointerFromString(fmt.Sprintf("%s/%s", r.pointer.String(), "sha1?")),
				Value: strings.TrimPrefix(r.Compiled.Sha1, "sha1:"),
			},
			patch.ReplaceOp{
				Path:  patch.MustNewPointerFromString(fmt.Sprintf("%s/%s", r.pointer.String(), "stemcell?/os")),
				Value: r.Stemcell.OS,
			},
			patch.ReplaceOp{
				Path:  patch.MustNewPointerFromString(fmt.Sprintf("%s/%s", r.pointer.String(), "stemcell?/version")),
				Value: r.Stemcell.Version,
			},
		)
	}

	return ops
}
