# Analysis

There are several built-in analyzers which can also be used directly, although they are typically used indirectly when using the `analysis` subcommands of specific artifact types (results from using these generic commands are not stored). Analysis results are JSON streams of objects.

Generate a raw analysis with the `generate` subcommand...

    $ boshua analysis generate --analyzer=releasemanifests.v1 openvpn-4.2.1.tgz | jq -c keys | head -n1
    ["parsed","path","raw"]

The `formatter` subcommand provides direct access to built-in formatters which can parse analysis results into more meaningful output...

    $ boshua analysis formatter stemcellimagefiles.v1 -h
    Available commands:
      ls         Show an ls-style list of files
      sha1sum    Show sha1 checksums (aliases: shasum)
      sha256sum  Show sha256 checksums
      sha512sum  Show sha512 checksums

Pass the results of a previous analysis via `STDIN` to use it...

    $ boshua analysis formatter releasemanifests.v1 properties --job=openvpn < releasemanifests.v1.jsonl

The previous command is equivalent to the following examples due to how the `results` subcommand of individual artifacts works...

    $ boshua release --release=openvpn/5.1.0 analysis --analyzer=releasemanifests.v1 results -- properties --job=openvpn
    $ boshua release --release=openvpn/5.1.0 analysis --analyzer=releasemanifests.v1 results --raw | boshua analysis formatter releasemanifests.v1 properties --job=openvpn


## Analyzers

The following default analyzers are included...

 * `releaseartifactfiles.v1` - extract file stats and checksums from a release tarball and its embedded tarballs
 * `releasemanifests.v1` - extract `release.MF` and job `spec` data from a release tarball
 * `stemcellimagefiles.v1` - extract file stats and checksums from the embedded image of a stemcell tarball
 * `stemcellmanifest.v1` - extract `stemcell.MF` data from a stemcell tarball
 * `stemcellpackages.v1` - extract and parse `packages.txt` data from a stemcell tarball

