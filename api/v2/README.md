## Find a Compiled Release

    > POST /lookup
    > Accept: application/json
    >
    > {
    >   "data": {
    >     "release": {
    >       "name": "openvpn",
    >       "version": "4.2.1",
    >       "checksum": {
    >         "type": "sha1",
    >         "value": "cbaca9fe8ceffb13ea1d481ad49fd5a706afbe9c"
    >       }
    >     },
    >     "stemcell": {
    >       "os": "ubuntu-trusty",
    >       "version": "3468.22"
    >     }
    >   }
    > }
    =
    < HTTP/1.0 200 OK
    < Content-Type: application/json
    <
    < {
    <   "compiled_release_version": {
    <     "url": "https://s3-external-1.amazonaws.com/dpb587-bosh-release-openvpn-us-east-1/compiled_releases/openvpn/openvpn-4.2.1-on-ubuntu-trusty-stemcell-3468.22-compiled-1.20180213074109.0.tgz",
    <     "checksums": [
    <       {
    <         "type": "sha1",
    <         "value": "3c313c203593572c5bf8237fb6bd32eebf33baba"
    <       },
    <       {
    <         "type": "sha256",
    <         "value": "f5465c2aaec28d5a44c560c2a87bfa50fe2db155f5a4f4ad492ab2af0f1599a9"
    <       }
    <     },
    <     "release": {
    <       "name": "openvpn",
    <       "version": "4.2.1",
    <       "checksum": {
    <         "type": "sha1",
    <         "value": "cbaca9fe8ceffb13ea1d481ad49fd5a706afbe9c"
    <       }
    <     },
    <     "stemcell": {
    <       "os": "ubuntu-trusty",
    <       "version": "3468.27"
    <     }
    <   }
    < }

## Schedule a Compiled Release

    > POST /schedule
    > Accept: application/json
    >
    > {
    >   "data": {
    >     "release": {
    >       "name": "openvpn",
    >       "version": "4.2.1",
    >       "checksum": {
    >         "type": "sha1",
    >         "value": "cbaca9fe8ceffb13ea1d481ad49fd5a706afbe9c"
    >       }
    >     },
    >     "stemcell": {
    >       "os": "ubuntu-trusty",
    >       "version": "3468.22"
    >     }
    >   }
    > }
    =
    < HTTP/1.0 202 Accepted
    < Content-Type: application/json
    <
    < {
    <   "status": "scheduled"
    < }
