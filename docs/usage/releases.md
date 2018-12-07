# Releases

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


## Compilations

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

