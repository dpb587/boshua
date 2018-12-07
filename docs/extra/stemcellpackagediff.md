# stemcellpackagediff

The [`stemcellpackagediff`](../../main/stemcellpackagediff) is a standalone CLI which uses the `stemcellpackages.v1` analyzer to compare packages between two versions of a stemcell.

    stemcellpackagediff ubuntu-xenial 97.12 97.15
    ~ intel-microcode (3.20180807a.0ubuntu0.16.04.1; was 3.20180425.1~ubuntu0.16.04.2)
    ~ linux-generic-hwe-16.04-edge (4.15.0.33.54; was 4.15.0.32.53)
    - linux-headers-4.15.0-32 (4.15.0-32.35~16.04.1)
    - linux-headers-4.15.0-32-generic (4.15.0-32.35~16.04.1)
    + linux-headers-4.15.0-33 (4.15.0-33.36~16.04.1)
    + linux-headers-4.15.0-33-generic (4.15.0-33.36~16.04.1)
    ~ linux-headers-generic-hwe-16.04-edge (4.15.0.33.54; was 4.15.0.32.53)
    - linux-image-4.15.0-32-generic (4.15.0-32.35~16.04.1)
    + linux-image-4.15.0-33-generic (4.15.0-33.36~16.04.1)
    ~ linux-image-generic-hwe-16.04-edge (4.15.0.33.54; was 4.15.0.32.53)
    ~ linux-libc-dev:amd64 (4.4.0-134.160; was 4.4.0-133.159)
    - linux-modules-4.15.0-32-generic (4.15.0-32.35~16.04.1)
    + linux-modules-4.15.0-33-generic (4.15.0-33.36~16.04.1)
    - linux-modules-extra-4.15.0-32-generic (4.15.0-32.35~16.04.1)
    + linux-modules-extra-4.15.0-33-generic (4.15.0-33.36~16.04.1)

It supports a `--format` argument where `markdown` (below) or `json` can be used to customize the output.

| ubuntu-xenial | 97.12 | 97.15 |
|:------------- | -----:| -----:|
| intel-microcode | 3.20180425.1~ubuntu0.16.04.2 | 3.20180807a.0ubuntu0.16.04.1 |
| linux-generic-hwe-16.04-edge | 4.15.0.32.53 | 4.15.0.33.54 |
| linux-headers-4.15.0-32 | 4.15.0-32.35~16.04.1 | &ndash; |
| linux-headers-4.15.0-32-generic | 4.15.0-32.35~16.04.1 | &ndash; |
| linux-headers-4.15.0-33 | &ndash; | 4.15.0-33.36~16.04.1 |
| linux-headers-4.15.0-33-generic | &ndash; | 4.15.0-33.36~16.04.1 |
| linux-headers-generic-hwe-16.04-edge | 4.15.0.32.53 | 4.15.0.33.54 |
| linux-image-4.15.0-32-generic | 4.15.0-32.35~16.04.1 | &ndash; |
| linux-image-4.15.0-33-generic | &ndash; | 4.15.0-33.36~16.04.1 |
| linux-image-generic-hwe-16.04-edge | 4.15.0.32.53 | 4.15.0.33.54 |
| linux-libc-dev:amd64 | 4.4.0-133.159 | 4.4.0-134.160 |
| linux-modules-4.15.0-32-generic | 4.15.0-32.35~16.04.1 | &ndash; |
| linux-modules-4.15.0-33-generic | &ndash; | 4.15.0-33.36~16.04.1 |
| linux-modules-extra-4.15.0-32-generic | 4.15.0-32.35~16.04.1 | &ndash; |
| linux-modules-extra-4.15.0-33-generic | &ndash; | 4.15.0-33.36~16.04.1 |

