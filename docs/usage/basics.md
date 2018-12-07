# Basic Usage

The primary tool is the [`boshua`](../../main/boshua) CLI. Generally, `-h` should be used to learn more about specific commands and available options.

The first level of commands are primarily geared towards specific artifacts (e.g. `release`, `stemcell`)...

    $ boshua -h
    Available commands:
      release     For working with releases
      stemcell    For working with stemcells

When working with a specific artifact type, several flags can be used for finding specific versions of the artifact...

    $ boshua release -h
    [release command options]
          --release=          The release in name/version format
          --release-name=     The release name
          --release-version=  The release version
          --release-checksum= The release checksum
          --release-url=      The release source URL
          --release-label=    The label(s) to filter releases by

Further subcommands can be used which deal with a specific artifact type...

    $ boshua release --release=openvpn/5.1.0 -h
    Available commands:
      analysis        For analyzing the release artifact
      analyzers       For showing the supported analyzers
      artifact        For showing the release artifact
      compilation     For working with compiled releases
      datastore       For interacting with release datastores
      download        For downloading the release locally
      upload-release  For uploading the release to BOSH

Most notably, these subcommands are common across most artifacts: `analysis`, `analyzers`, `artifact`, `datastore`, and `download`.

The `analysis` subcommands have some shared usage as well...

    $ boshua release --release=openvpn/5.1.0 analysis -h
    [analysis command options]
          --analyzer=         The analyzer name

    Available commands:
      artifact       For showing the analysis artifact
      download       For downloading the analysis locally
      results        For showing the results of an analysis
