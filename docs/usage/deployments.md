# Deployments

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

By default, if a compilation is not already available and it can be scheduled, this command will queue and block until the shared compilation artifact becomes available.

