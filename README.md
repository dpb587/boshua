# boshua

For providing, using, and inspecting artifacts of [BOSH](https://bosh.io/).


## Usage



See the following for some specific examples of usage.


### Deployment Manifests

Convert a manifest referencing release sources to compiled releases...

    $ bosh deployment use-compiled-releases < manifest.yml
    TODO sample


### Releases

Showing the tarball of a release...

    $ boshua release --release=openvpn/5.0.0 --release-checksum=0b08f569dc18b042845897a0490d541f96f96951
    file    openvpn-5.0.0.tgz
    url     https://s3-external-1.amazonaws.com/bosh-hub-release-tarballs/7f98eb62-f111-461f-71a1-70853052d90c
    sha1    0b08f569dc18b042845897a0490d541f96f96951
    ...


### Stemcells

Showing the tarball of a stemcell...

    $ boshua stemcell --stemcell=bosh-aws-xen-hvm-ubuntu-trusty-go_agent/3541.12 --light-stemcell
    file light-bosh-stemcell-3541.12-aws-xen-hvm-ubuntu-trusty-go_agent.tgz
    url  https://s3.amazonaws.com/bosh-aws-light-stemcells/light-bosh-stemcell-3541.12-aws-xen-hvm-ubuntu-trusty-go_agent.tgz
    sha1 e2f9840e7ed3eb2ccdf4c39f3a7b49e35e1ad8ec
    ...

Show the filesystem of a stemcell...

    $ boshua stemcell --stemcell=bosh-aws-xen-hvm-ubuntu-trusty-go_agent/3541.12 analysis --analyzer=stemcellimagefiles.v1
    drwxr-xr-x 2 root root       0 Apr  6 18:43 /bin
    -rwxr-xr-x 1 root root 1021112 May 16  2017 /bin/bash
    -rwxr-xr-x 3 root root   31152 Oct 21  2013 /bin/bunzip2
    -rwxr-xr-x 3 root root   31152 Oct 21  2013 /bin/bzcat
    lrwxrwxrwx 1 root root       6 Oct 21  2013 /bin/bzcmp -> bzdiff
    -rwxr-xr-x 1 root root    2140 Oct 21  2013 /bin/bzdiff
    lrwxrwxrwx 1 root root       6 Oct 21  2013 /bin/bzegrep -> bzgrep
    -rwxr-xr-x 1 root root    4877 Oct 21  2013 /bin/bzexe
    lrwxrwxrwx 1 root root       6 Oct 21  2013 /bin/bzfgrep -> bzgrep
    ...

Show the packages of a stemcell...

    $ boshua stemcell --stemcell=bosh-aws-xen-hvm-ubuntu-trusty-go_agent/3541.12 analysis --analyzer=stemcellpackages.v1
    Desired=Unknown/Install/Remove/Purge/Hold
    | Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend
    |/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)
    ||/ Name                                Version                                    Architecture Description
    +++-===================================-==========================================-============-===============================================================================
    ii  adduser                             3.113+nmu3ubuntu3                          all          add and remove users and groups
    ii  anacron                             2.3-20ubuntu1                              amd64        cron-like program that doesn't go by time
    ii  apparmor                            2.10.95-0ubuntu2.6~14.04.3                 amd64        user-space parser utility for AppArmor
    ii  apparmor-utils                      2.10.95-0ubuntu2.6~14.04.3                 amd64        utilities for controlling AppArmor
    ...


## Limitations

 * TODO security
 * when patching deployment manifests to use compiled releases...
    * releases must already specify expected tarball checksums
    * explicit versions (not `latest`) must be used for `releases` and `stemcells`
    * multiple stemcells must not be used


## Futures

 * standalone compilations
 * smarter factories for knowing writeable indices


## License

[MIT License](LICENSE)
