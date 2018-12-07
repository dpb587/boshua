# Configuration

The default configuration lives in `~/.config/boshua/config.yml` - a YAML file describing the various datastores that `boshua` can reference. See the partial documentation in [`config/config.go`](config/config.go) or examples in [`examples`](examples).

When a remote server is available, `$BOSHUA_SERVER` or `--default-server` may be used to automatically configure remote lookups and avoid the need for a configuration file. When building [`boshua`](main/boshua), a default server can be embedded with [link options](https://golang.org/cmd/link/) and the `main.defaultServer` variable.


## Stemcell Providers

 * [`boshioindex`](stemcellversion/datastore/boshioindex) - [bosh-io/stemcells-core-index](https://github.com/bosh-io/stemcells-core-index)-style
 * [`boshua.v2`](stemcellversion/datastore/boshua.v2) - query a remote boshua API server


## Release Providers

 * [`boshioindex`](releaseversion/datastore/boshioindex) - [bosh-io/releases-index](https://github.com/bosh-io/releases-index)-style
 * [`boshreleasedir`](releaseversion/datastore/boshreleasedir) - directly reference a release repository (i.e. any BOSH repository)
 * [`boshua.v2`](releaseversion/datastore/boshua.v2) - query a remote boshua API server
 * [`metalinkrepository`](releaseversion/datastore/metalinkrepository) - refer to a metalink repository of pre-built release tarballs (e.g. [dpb587/openvpn-bosh-release](https://github.com/dpb587/openvpn-bosh-release/tree/artifacts/release/stable))
 * [`trustedtarball`](releaseversion/datastore/trustedtarball) - dynamically generate artifacts for queried tarballs (e.g. internal dev builds)


## Release Compilation Providers

 * [`boshua.v2`](releaseversion/compilation/datastore/boshua.v2) - query a remote boshua API server
 * [`contextualosmetalinkrepository`](releaseversion/compilation/datastore/contextualosmetalinkrepository) - refer to a metalink repository, segmented by OS name and version (e.g. [dpb587/openvpn-bosh-release](https://github.com/dpb587/openvpn-bosh-release/tree/artifacts/compiled-release/stable))
 * [`contextualrepoosmetalinkrepository`](releaseversion/compilation/datastore/contextualrepoosmetalinkrepository) - refer to a metalink repository, segmented by `repo`-label and OS name and version (e.g. internal, shared compilation repository)


## Scheduler Providers

 * [`boshua.v2`](task/scheduler/boshua.v2) - schedule via a remote boshua API server
 * [`concourse`](task/scheduler/concourse) - schedule tasks via concourse pipelines
 * [`localexec`](task/scheduler/localexec) - run tasks locally


## Download Handlers


### URLs

 * [`s3`](artifact/downloader/url/s3) - authenticated downloads from S3 buckets
 * [`pivnet`](artifact/downloader/url/pivnet) - authenticated downloads from [Pivotal Network](https://network.pivotal.io/)
 * `http` - download from http(s) servers; default
 * `ftp` - download from ftp servers; default


### MetaURLs

 * `boshreleasesource` - build release tarballs directly from repositories; default

