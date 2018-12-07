# Stemcells

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

