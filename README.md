# boshua

For providing, using, and inspecting artifacts of [BOSH](https://bosh.io/).

> bosh unofficial artifacts


## Usage

See the following for some specific examples of usage.


### Releases

Showing the tarball of a release...

    $ boshua release --release=openvpn/5.0.0
    file	openvpn-5.0.0.tgz
    url	https://s3-external-1.amazonaws.com/dpb587-bosh-release-openvpn-us-east-1/artifacts/release/1b77cbd51a3debefcb06f2ad5311d872f056dbe9
    sha1	1b77cbd51a3debefcb06f2ad5311d872f056dbe9
    ...

Getting the manifest files of a release...

    $ boshua release --release=openvpn/5.0.0 analysis --analyzer=releasemanifests.v1 results --raw

Getting the properties for a job of a release...

    $ boshua release --release=openvpn/5.0.0 analysis --analyzer=releasemanifests.v1 results -- properties --job=openvpn
    server     VPN IP and netmask (basis of the IP pool which the server will allocate to clients)
    tls_cipher A colon-separated list of allowable TLS ciphers
    tls_crl    Certificate Revocation List (`X509 CRL`, including the begin/end markers)
    dh_pem     Diffie-Hellmann Key (`DH PARAMETERS`, including the begin/end markers)
    ...


#### Compilations

Finding the compilation of a release...

    $ boshua release --release=openvpn/4.1.0 compilation --os=ubuntu-trusty/3468.13
    file	openvpn-4.1.0-on-ubuntu-trusty-stemcell-3468.13.tgz
    url	https://s3-external-1.amazonaws.com/dpb587-bosh-release-openvpn-us-east-1/compiled_releases/openvpn/openvpn-4.1.0-on-ubuntu-trusty-stemcell-3468.13-compiled-1.20171209113453.0.tgz
    sha1	d278e2a37c486beabd0a9ffd2426e58b38172842
    ...

Showing the checksums of the files of a compiled release...

    $ boshua release --release=openvpn/4.1.0 compilation --os=ubuntu-trusty/3468.12 analysis --analyzer=releaseartifactfiles.v1 results -- sha1sum
    fd01d7b1fa7929906db7486943e3c68510794d01  compiled_packages/openvpn.tgz!external/openssl/include/openssl/aes.h
    3df405ed1c9876d25551b732f9c2985ecbf1fdf6  compiled_packages/openvpn.tgz!external/openssl/include/openssl/asn1.h
    9fe8dd066ed9109c09862222a25b15bf109ad34c  compiled_packages/openvpn.tgz!external/openssl/include/openssl/asn1_mac.h
    4642be4516e5af219da061da7b4edfa948bd590e  compiled_packages/openvpn.tgz!external/openssl/include/openssl/asn1t.h
    ...


### Stemcells

Showing the tarball of a stemcell...

    $ boshua stemcell --stemcell=bosh-aws-xen-hvm-ubuntu-trusty-go_agent/3541.12 --stemcell-flavor=light
    file	light-bosh-stemcell-3541.12-aws-xen-hvm-ubuntu-trusty-go_agent.tgz
    url	https://s3.amazonaws.com/bosh-aws-light-stemcells/light-bosh-stemcell-3541.12-aws-xen-hvm-ubuntu-trusty-go_agent.tgz
    sha1	e2f9840e7ed3eb2ccdf4c39f3a7b49e35e1ad8ec
    ...

Show the filesystem of a stemcell...

    $ boshua stemcell --stemcell=bosh-aws-xen-hvm-ubuntu-trusty-go_agent/3541.12 analysis --analyzer=stemcellimagefiles.v1 results -- ls
    drwxr-xr-x - root root       0 Apr  6 18:43 /bin
    -rwxr-xr-x - root root 1021112 May 16  2017 /bin/bash
    -rwxr-xr-x - root root   31152 Oct 21  2013 /bin/bunzip2
    -rwxr-xr-x - root root   31152 Oct 21  2013 /bin/bzcat
    lrwxrwxrwx - root root       6 Oct 21  2013 /bin/bzcmp -> bzdiff
    -rwxr-xr-x - root root    2140 Oct 21  2013 /bin/bzdiff
    lrwxrwxrwx - root root       6 Oct 21  2013 /bin/bzegrep -> bzgrep
    -rwxr-xr-x - root root    4877 Oct 21  2013 /bin/bzexe
    lrwxrwxrwx - root root       6 Oct 21  2013 /bin/bzfgrep -> bzgrep
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


### Deployment Manifests

*TODO resurrecting post-refactor*

Convert a manifest referencing release sources to compiled releases...

    $ bosh deployment use-compiled-releases < manifest.yml


### GraphQL

 > http://localhost:4508/api/v2/graphql?query={...}

*TODO experimental/planning*

Fetch tarball of a compiled release...

    {
      release(name: "openvpn", version: "5.1.0") {
        compilation(os: "ubuntu-trusty", version: "3586.12") {
          tarball {
            hash(type: "sha1"),
            url
          }
        }
      }
    }

Fetch jobs and packages of a release (this would need to error if the analysis had not been performed? better to recommend analysis?)...

    {
      release(name: "openvpn", version: "5.1.0") {
        jobs {
          name,
          consumes,
          provides,
          dependencies,
          properties
        }

        packages {
          name,
          dependencies
        }
      }
    }

Fetch analysis artifacts of a release...

    {
      release(name: "openvpn", version: "5.1.0") {
        analysis {
          results(analyzers: ["releaseartifactfiles.v1", "releasemanifests.v1"]) {
            analyzer
            status
            artifact {
              hash(type: "sha1"),
              url
            }
          }
        }
      }
    }

Fetch specialized analysis results...

    {
      release(name: "openvpn", version: "5.1.0") {
        analysis {
          releaseartifactfilesV1 {
            results {
              totalCount
              pageInfo {
                endCursor
                hasNextPage
              }
              nodes {
                artifact
                result {
                  type
                  path
                  link
                  size
                  mode
                }
              }
            }
          }
        }
      }
    }

Request an analysis be made...

    mutation {
      executeAnalysis(
        analyzer: "releaseartifactfiles.v1",
        subject: {
          release: {
            name: "openvpn",
            version: "5.1.0",
            checksum: "sha1:b42eb85e5f074c26b065956cc9b8a6d69208f8a0"
          }
        }
      ) {
        status
      }
    }

Lookup a stemcell...

    {
      stemcell(os: "ubuntu-xenial", version: "87.3", iaas: "aws") {
        url
        sha1
      }
    }


### Web UI

> http://localhost:4508/webui/

*TODO experimental/playing*

 * [releases.html](http://localhost:4508/webui/releases.html)
 * [stemcells.html](http://localhost:4508/webui/stemcells.html)


## Concepts

 * **Artifact** - an artifact represents something of interest and generally refers to a permanent blob of data somewhere, such as a BOSH release tarball stored on Amazon S3. Artifacts are usually identified by a couple pieces of information (e.g. name, version, checksum). There are several primary types of artifacts, each represented by a top-level CLI command.
    * stemcell - a particular version of a BOSH stemcell for a given IaaS
    * release - a particular version of a BOSH release
    * compiled-release - a particular version of a BOSH release that has been compiled against a particular OS and version
 * **Analysis** - generated metadata about an artifact. There are several different analyzers, all of which generate JSON data. Most analyzers have default formatters for rendering the data in a meaningful way.
 * **Labels** - used to label/tag artifacts for logical categorization. Recommended to use path-style for hierarchical taxonomies. Examples...
    * `repo/*` - to identify source of the release; e.g. `repo/github.com/dpb587/openvpn-bosh-release`
    * `stability/(alpha|beta|rc|stable)` - identify stability of artifacts
    * `tag/*` - tag-based navigation; e.g. `tag/cpi`, `tag/networking`
    * `deprecated` - deprecated


## Limitations

 * when patching deployment manifests to use compiled releases...
    * releases must already specify expected tarball checksums
    * explicit versions (not `latest`) must be used for `releases` and `stemcells`
    * multiple stemcells must not be used


## Futures

 * clean. up.
 * namespacing git repository settings in config
 * mirror rewrites for proxying upstream artifacts
 * authentication?
 * logging


## License

[MIT License](LICENSE)
