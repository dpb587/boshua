# Concepts

First, let's define some of the key concepts that `boshua` is built on...


## Artifacts

An artifact represents something of interest to BOSH and generally refers to a permanent blob of data somewhere. They are typically identified by several properties such as name, version, and/or checksum. For example, a BOSH release tarball that lives in Amazon S3 and can be verified by its sha-1 checksum.


There are several primary types of artifacts, each represented by a top-level CLI command or API endpoint.


### Stemcell Artifacts

A stemcell represents a typical BOSH stemcell. It can be referenced with the following properties:

 * `os` - the published OS name (e.g. `ubuntu-xenial`, `windows2016`)
 * `version` - the version of the stemcell
 * `iaas` - the IaaS the stemcell was built for (e.g. `aws`, `google`)
 * `hypervisor` - the IaaS hypervisor the stemcell was built for (e.g. `xen`, `kvm`)
 * `flavor` - whether the stemcell is `heavy` or `light`


### Releases Artifacts

A release represents a typical BOSH release. It can be referenced with the following properties:

 * `name` - the BOSH release name, as it would be found in a deployment manifest or in `config/final.yml`
 * `version` - the BOSH release version, as it would be found in a deployment manifest or in `releases/*/*.yml`
 * `checksum` - for a pre-existing release tarball, one of its pre-calculated checksums (i.e. `sha1`/`sha256`)
 * `url` - for a pre-existing release tarball, one of its download URIs

In addition to release artifacts, there are also compilation artifacts. These artifacts represent a specific compilation of a release artifact against a stemcell. They are the cross-reference between a release artifact and stemcell artifact.


### Labels Metadata

In addition to individual artifact properties, artifacts may also be classified by arbitrary labels which may be more meaningful to consumers. For hierarchical taxonomies, path-style values with forward slashes should be used. Standard label conventions are:

 * `repo/*` - to identify the source of the release; e.g. `repo/github.com/dpb587/openvpn-bosh-release`
 * `stability/(alpha|beta|rc|stable)` - identify stability of artifacts
 * `tag/*` - general, tag-based navigation; e.g. `tag/cpi`, `tag/networking`
 * `deprecated` - to identify artifacts which are deprecated


## Analyses

The artifacts are independently useful, but there is often much more information which can be derived from it (e.g. what operating system packages are included in a particular stemcell version). Analysis results are generated metadata which is associated with a particular artifact.

Analyzers are used to generate specific types of raw analysis results about an artifact. Each artifact type has several builtin analyzers. For example, an analyzer might audit and report on the expected files and checksums in a release tarball.

Formatters are tools for interpreting the raw results and providing them in a more meaningful way. Most analyzers have several builtin formatters.


## Datastores

A datastore is something which can find and/or store details about artifacts in a permanent way (e.g. a BOSH release repository having release information). Each artifact type has several supported datastores, and datastores can delegate to other, possibly remote datastores (e.g. through APIs).


## Schedulers

A scheduler is used for executing work when results are not already available (e.g. compiling a release). Several types of schedulers are supported to support running work locally or in Docker, remotely on Concourse, or remotely through an API.
