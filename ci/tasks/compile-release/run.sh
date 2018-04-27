#!/bin/bash

set -eu -o pipefail

tar -xzf release/*.tgz $( tar -tzf release/*.tgz | grep release.MF$ )
release_name=$( grep '^name:' release.MF | awk '{print $2}' | tr -d "\"'" )
release_version=$( grep '^version:' release.MF | awk '{print $2}' | tr -d "\"'" )

tar -xzf stemcell/*.tgz stemcell.MF
stemcell_name=$( grep '^name:' stemcell.MF | awk '{print $2}' | tr -d "\"'" )
stemcell_os=$( grep '^operating_system:' stemcell.MF | awk '{print $2}' | tr -d "\"'" )
stemcell_version=$( grep '^version:' stemcell.MF | awk '{print $2}' | tr -d "\"'" )

export BOSH_DEPLOYMENT=compilation

cat > deployment.yml <<EOF
name: "$BOSH_DEPLOYMENT"
instance_groups: []
releases:
- name: "$release_name"
  version: "$release_version"
stemcells:
- alias: "default"
  name: "$stemcell_name"
  version: "$stemcell_version"
update:
  canaries: 1
  canary_watch_time: 1
  max_in_flight: 1
  update_watch_time: 1
EOF

bosh-director/bosh -n deploy deployment.yml

bosh-director/bosh export-release "$release_name/$release_version" "$stemcell_os/$stemcell_version"

bosh-director/bosh task --event 4 > compiled-release/events.json

version=$( 1.$( date -u +%Y%m%d%H%M%S ).0 )

echo -n "$version" > compiled-release/version

mv *.tgz compiled-release/$release_name-$release_version-on-$stemcell_os-stemcell-$stemcell_version-compiled-$version.tgz

bosh-director/bosh inspect-release "$release_name/$release_version"

bosh-director/bosh -n delete-deployment
