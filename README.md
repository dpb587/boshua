# boshua

For providing, using, and inspecting artifacts of [BOSH](https://bosh.io/).

> bosh unofficial artifacts


## Core Concepts

First, let's define some of the terminology this uses...

 * **Artifact** - an artifact represents something of interest and generally refers to a permanent blob of data somewhere (e.g. a BOSH release tarball stored on Amazon S3). Artifacts are usually identified by a couple pieces of canonical information (e.g. name, version, URI, checksum). There are several primary types of artifacts, each represented by a top-level CLI command or API endpoint.
    * **Stemcell** - a particular version of a BOSH stemcell for a given IaaS
    * **Release** - a particular version of a BOSH release
       * **Compilation** - a particular version of a BOSH release that has been compiled against a particular stemcell or OS and version
    * **Labels** - used to label artifacts for logical categorization. For hierarchical taxonomies, path-style values with forward slashes can be used. Some standardized label conventions are...
       * `repo/*` - to identify the source of the release; e.g. `repo/github.com/dpb587/openvpn-bosh-release`
       * `stability/(alpha|beta|rc|stable)` - identify stability of artifacts
       * `tag/*` - general, tag-based navigation; e.g. `tag/cpi`, `tag/networking`
       * `deprecated` - to identify artifacts which are deprecated
 * **Analysis** - an artifact is independently useful, but there is often much more information which can be derived from it (e.g. what OS packages are included in a particular stemcell version). Analysis results are generated metadata which is affiliated with a particular artifact.
    * **Analyzer** - analyzers are used to generate specific types of metadata about an artifact. Each artifact type has several builtin analyzers.
    * **Formatters** - formatters are built-in tools for interpreting the raw results and providing them in a more meaningful way.
 * **Datastore** - a datastore is something which can find and/or store details about artifacts and analysis in a permanent way (e.g. a BOSH release repository having release information). Each artifact type has several supported datastores, and datastores can delegate to other, possibly remote datastores (e.g. through APIs).
 * **Scheduler** - a scheduler is used for executing work when the results are not already available (e.g. compiling a release). Several types of schedulers are supported to support running work locally or in Docker, remotely on Concourse, or remotely through an API.


## Usage


### CLI

See the following for some specific examples of usage.


#### Releases

Finding the tarball of a release...

    $ boshua release --release=openvpn/5.0.0
    file   openvpn-5.0.0.tgz
    url    https://s3-external-1.amazonaws.com/dpb587-bosh-release-openvpn-us-east-1/artifacts/release/1b77cbd51a3debefcb06f2ad5311d872f056dbe9
    sha1   1b77cbd51a3debefcb06f2ad5311d872f056dbe9
    sha256 02965881f86b36b311768b154dadbef4522a8cccb81e1b6531c7db05848869aa

Showing the `release.MF` data of a release...

    $ boshua release --release=openvpn/5.0.0 analysis --analyzer=releasemanifests.v1 results -- spec --release
    name: openvpn
    version: 5.0.0
    commit_hash: 0f8966c
    uncommitted_changes: false
    ...

Showing the properties for a job of a release...

    $ boshua release --release=openvpn/5.0.0 analysis --analyzer=releasemanifests.v1 results -- properties --job=openvpn
    server     VPN IP and netmask (basis of the IP pool which the server will allocate to clients)
    tls_cipher A colon-separated list of allowable TLS ciphers
    tls_crl    Certificate Revocation List (`X509 CRL`, including the begin/end markers)
    dh_pem     Diffie-Hellmann Key (`DH PARAMETERS`, including the begin/end markers)
    ...


##### Compilations

Getting the compilation of a release on a stemcell...

    $ boshua release --release=openvpn/4.1.0 compilation --os=ubuntu-trusty/3468.13
    file   openvpn-4.1.0-on-ubuntu-trusty-stemcell-3468.13.tgz
    url    https://s3-external-1.amazonaws.com/dpb587-bosh-release-openvpn-us-east-1/compiled_releases/openvpn/openvpn-4.1.0-on-ubuntu-trusty-stemcell-3468.13-compiled-1.20171209113453.0.tgz
    sha1   d278e2a37c486beabd0a9ffd2426e58b38172842
    sha256 02120c9f1d084e232c0a996f7fa54e0e41c8b53c72cdb1003085108311929362

Uploading the compilation to the director (or showing the command to)...

    $ boshua release --release=openvpn/4.1.0 compilation --os=ubuntu-trusty/3468.13 upload-release --cmd
    bosh upload-release --name=openvpn --version=4.1.0 \
      https://s3-external-1.amazonaws.com/dpb587-bosh-release-openvpn-us-east-1/compiled_releases/openvpn/openvpn-4.1.0-on-ubuntu-trusty-stemcell-3468.13-compiled-1.20171209113453.0.tgz \
      --sha1=d278e2a37c486beabd0a9ffd2426e58b38172842 \
      --stemcell=ubuntu-trusty/3468.13

Getting an ops file for using the compiled release in a manifest...

    $ boshua release --release=openvpn/4.1.0 compilation --os=ubuntu-trusty/3468.13 ops-file
    - path: /releases/name=openvpn?
      type: replace
      value:
        name: openvpn
        sha1: md5:9cc79bee6180ef5e9f9b96606bf252bd
        stemcell:
          os: ubuntu-trusty
          version: "3468.13"
        url: https://s3-external-1.amazonaws.com/dpb587-bosh-release-openvpn-us-east-1/compiled_releases/openvpn/openvpn-4.1.0-on-ubuntu-trusty-stemcell-3468.13-compiled-1.20171209113453.0.tgz
        version: 4.1.0

Showing the checksums of the files of a compiled release...

    $ boshua release --release=openvpn/4.1.0 compilation --os=ubuntu-trusty/3468.12 analysis --analyzer=releaseartifactfiles.v1 results -- sha1sum
    ...
    7edc92307679f49446037387effa6c642c05e2e0  compiled_packages/openvpn.tgz!share/doc/openvpn/COPYRIGHT.GPL
    67766b2d0c67c36841e77c6b05673a702559371b  compiled_packages/openvpn.tgz!share/doc/openvpn/COPYING
    99e42912c49c8cd676000c00f2dd51c1795cb4f4  compiled_packages/openvpn.tgz!share/man/man8/openvpn.8
    e0ebceb7f4f638aca7210001c828d6f889a8128f  compiled_packages/openvpn.tgz!lib/openvpn/plugins/openvpn-plugin-down-root.so
    6eb2e481af90d6060a61a889a8641dc1e5e75331  compiled_packages/openvpn.tgz!lib/openvpn/plugins/openvpn-plugin-down-root.la


#### Stemcells

Finding the tarball of a stemcell...

    $ boshua stemcell --stemcell=bosh-aws-xen-hvm-ubuntu-trusty-go_agent/3541.12 --stemcell-flavor=light
    file  light-bosh-stemcell-3541.12-aws-xen-hvm-ubuntu-trusty-go_agent.tgz
    url   https://s3.amazonaws.com/bosh-aws-light-stemcells/light-bosh-stemcell-3541.12-aws-xen-hvm-ubuntu-trusty-go_agent.tgz
    sha1  e2f9840e7ed3eb2ccdf4c39f3a7b49e35e1ad8ec
    sha256 23884d534e4f5f946234ff3caf4240f20a37473b6afa0fcb5ba0f5bca3f9de3c

Show the filesystem of a stemcell...

    $ boshua stemcell --stemcell=bosh-aws-xen-hvm-ubuntu-trusty-go_agent/3541.12 --stemcell-flavor=light analysis --analyzer=stemcellimagefiles.v1 results -- ls
    drwxr-xr-x - root root       0 Apr  6 18:43 /bin
    -rwxr-xr-x - root root 1021112 May 16  2017 /bin/bash
    -rwxr-xr-x - root root   31152 Oct 21  2013 /bin/bunzip2
    -rwxr-xr-x - root root   31152 Oct 21  2013 /bin/bzcat
    lrwxrwxrwx - root root       6 Oct 21  2013 /bin/bzcmp -> bzdiff
    ...

Show the packages of a stemcell...

    $ boshua stemcell --stemcell=bosh-aws-xen-hvm-ubuntu-trusty-go_agent/3541.12 --stemcell-flavor=light analysis --analyzer=stemcellpackages.v1 results -- contents
    Desired=Unknown/Install/Remove/Purge/Hold
    | Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend
    |/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)
    ||/ Name                                Version                                    Architecture Description
    +++-===================================-==========================================-============-===============================================================================
    ii  adduser                             3.113+nmu3ubuntu3                          all          add and remove users and groups
    ii  anacron                             2.3-20ubuntu1                              amd64        cron-like program that doesn't go by time
    ii  apparmor                            2.10.95-0ubuntu2.6~14.04.3                 amd64        user-space parser utility for AppArmor
    ii  apparmor-utils                      2.10.95-0ubuntu2.6~14.04.3                 amd64        utilities for controlling AppArmor
    ii  apt                                 1.0.1ubuntu2.17                            amd64        commandline package manager
    ...


#### Deployment Manifests

Convert a manifest referencing release sources to compiled releases...

    $ bosh deploy <( boshua deployment use-compiled-releases < manifest.yml )
      - name: openvpn
        version: 5.1.0
    -   sha1: b42eb85e5f074c26b065956cc9b8a6d69208f8a0
    +   sha1: sha512:ea3c1185076d52b87064951e91dd8885ca62f045dd4e1a17305e6a90a1901cb8d89ea097773e232bdbca2455be746672ea7be93371915597c574cb6933b7c13b
    -   url: https://s3-external-1.amazonaws.com/dpb587-bosh-release-openvpn-us-east-1/artifacts/release/b42eb85e5f074c26b065956cc9b8a6d69208f8a0
    +   url: https://s3-external-1.amazonaws.com/dpb587-test-20140414a-us-east-1/compiled-release/2f/474fe4787338086f4e0cb34207c3f687dabe16
    +   stemcell:
    +     os: ubuntu-trusty
    +     version: '3586.27'

Some caveats for automatically converting manifests...

 * explicit versions must be used for `releases` and `stemcells` (not `latest` or `x.latest`)
 * releases should specify canonical properties (e.g. absolute URLs or tarball checksums)
 * manifests with multiple stemcells are not supported


#### Server

The CLI provides an HTTP server to allow remote querying and execution of commands. By default, it will listen on `127.0.0.1:4508`.

    $ boshua server


##### CLI Downloads

The `/cli/` endpoint can be used for providing binaries for download.


##### Web UI

The `/ui/` endpoint can be used for hosting simple HTML pages.


##### GraphQL API

The `/api/v2/graphql` endpoint provides a GraphQL API with query and mutation support.

*This API has further changes pending; it is not stable.*


###### Query: `release`

Get information about a release...

    {
      release(name: String, version: String, url: String, checksum: String, labels: [String]) {
        name
        version
        labels

        tarball { #artifact }

        compilation(os: String, version: String) {
          tarball { #artifact }
          analysis { #analysis }
        }

        analysis {
          results(analyzers: [String]) {
            analyzer
            artifact { #artifact }
          }
        }
      }
    }


###### Query: `stemcell`

Get information about a stemcell...

    {
      stemcell(iaas: String, hypervisor: String, os: String, version: String, diskFormat: String, flavor: String, labels: [String]) {
        iaas
        hypervisor
        os
        flavor
        diskFormat
        version

        tarball { #artifact }
        analysis { #analysis }
      }
    }


    ###### Mutation: `scheduleCompilation`

    Schedule compilation for a release...

        mutation {
          scheduleReleaseCompilation(name: String, version: String, url: String, checksum: String, osName: String, osVersion: String) {
            status
          }
        }

###### Mutation: `scheduleStemcellAnalysis`

Schedule analysis of a stemcell...

    mutation {
      scheduleStemcellAnalysis( #stemcellParams, analyzer: String) {
        status
      }
    }

###### Mutation: `scheduleReleaseAnalysis`

Schedule analysis of a stemcell...

    mutation {
      scheduleReleaseAnalysis( #releaseParams, analyzer: String) {
        status
      }
    }

###### Mutation: `scheduleReleaseCompilationAnalysis`

Schedule analysis of a stemcell...

    mutation {
      scheduleReleaseAnalysis( #releaseParams, #releaseCompilationParams, analyzer: String) {
        status
      }
    }


### Library



### Configuration


## History

This project is the amalgamation of several former, smaller projects; but now focused on CLI+API, public+private deployability, and dynamic execution.

 * [bosh-stemcell-metadata-scripts](https://github.com/dpb587/bosh-stemcell-metadata-scripts) - repository of scripts for extracting package lists, filesystem details, and metadata from stemcells
 * [bosh-stemcell-metadata](https://github.com/dpb587/bosh-stemcell-metadata) - repository of pre-computed results from `bosh-stemcell-metadata-scripts` for specific stemcell lines
 * [bosh-release-compiler](https://github.com/dpb587/bosh-release-compiler) - simple repository of tasks for compiling releases with concourse
 * [bosh-compiled-releases](https://github.com/dpb587/bosh-compiled-releases) - repository for tracking compiled releases from shared environments or imported from external sources, and CLI for rewriting deployment manifests to use them


## License

[MIT License](LICENSE)
