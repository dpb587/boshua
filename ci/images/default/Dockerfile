FROM ubuntu:16.04
RUN true \
  && apt-get update \
  && apt-get install -y curl git openssl build-essential zlibc zlib1g-dev ruby ruby-dev libxslt1-dev libxml2-dev libssl-dev readline-common libreadline-dev libyaml-dev libsqlite3-dev sqlite3 multipath-tools \
  && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN true \
  && curl -fLo /usr/local/bin/bosh https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-3.0.1-linux-amd64 \
  && echo "ccc893bab8b219e9e4a628ed044ebca6c6de9ca0  /usr/local/bin/bosh" | sha1sum -c - \
  && chmod +x /usr/local/bin/bosh
RUN true \
  &&curl -fLo /usr/local/bin/meta4 https://github.com/dpb587/metalink/releases/download/v0.2.0/meta4-0.2.0-linux-amd64 \
  && echo "9b7afe863d95ed8d99886dbe1878619a71b97725  /usr/local/bin/meta4" | sha1sum -c - \
  && chmod +x /usr/local/bin/meta4
RUN true \
  && mkdir -p $HOME/.bosh/downloads \
  && curl -fLo $HOME/.bosh/downloads/0bee16bec9081702b63fe9e08345cf8051073a48-842e8dce585a4f04bd461667f77d468903e452ea https://bosh.io/d/stemcells/bosh-google-kvm-ubuntu-trusty-go_agent?v=3541.10 \
    && echo "842e8dce585a4f04bd461667f77d468903e452ea  $HOME/.bosh/downloads/0bee16bec9081702b63fe9e08345cf8051073a48-842e8dce585a4f04bd461667f77d468903e452ea" | sha1sum -c - \
  && curl -fLo $HOME/.bosh/downloads/57853bf6ff4474e15966faf3840e966ed73a8719-6ad87252c151f38d7e82cb93518829065ff27061 https://s3.amazonaws.com/bosh-compiled-release-tarballs/garden-linux-0.342.0-ubuntu-trusty-3541.10-20180327-211558-518479392-20180327211609.tgz?versionId=VOtVNvAlBiJ23VT5aLKlsKS.VoUoWsj0 \
    && echo "6ad87252c151f38d7e82cb93518829065ff27061  $HOME/.bosh/downloads/57853bf6ff4474e15966faf3840e966ed73a8719-6ad87252c151f38d7e82cb93518829065ff27061" | sha1sum -c - \
  && curl -fLo $HOME/.bosh/downloads/8126d6627029c3db5dafd63c2d4d147ab452d9f7-78d79f08ff5001cc2a24f572837c7a9c59a0e796 https://bosh.io/d/github.com/cloudfoundry/os-conf-release?v=18 \
    && echo "78d79f08ff5001cc2a24f572837c7a9c59a0e796  $HOME/.bosh/downloads/8126d6627029c3db5dafd63c2d4d147ab452d9f7-78d79f08ff5001cc2a24f572837c7a9c59a0e796" | sha1sum -c - \
  && curl -fLo $HOME/.bosh/downloads/8ced7086aa9155cccd57001fcefbd08b6e50d4ec-edaa548ac956d75ba416edfbd2621ef43ebb7f78 https://s3.amazonaws.com/bosh-compiled-release-tarballs/bosh-warden-cpi-40-ubuntu-trusty-3541.10-20180330-002031-436165203-20180330002036.tgz?versionId=0jdGpqM6RlJmgXErznZWErYG7UDaA1i2 \
    && echo "edaa548ac956d75ba416edfbd2621ef43ebb7f78  $HOME/.bosh/downloads/8ced7086aa9155cccd57001fcefbd08b6e50d4ec-edaa548ac956d75ba416edfbd2621ef43ebb7f78" | sha1sum -c - \
  && curl -fLo $HOME/.bosh/downloads/918bf014d10fc0f0ef61eb4fc6ffc93a205bd66c-a8f8f5138ef1e3b356bd5078a32694625e8846a5 https://s3.amazonaws.com/bosh-compiled-release-tarballs/bosh-265.2.0-ubuntu-trusty-3541.10-20180327-211954-148015419-20180327211957.tgz?versionId=cMRVRXQgYNbDeRdF52IyQ1yNBu7klA5j \
    && echo "a8f8f5138ef1e3b356bd5078a32694625e8846a5  $HOME/.bosh/downloads/918bf014d10fc0f0ef61eb4fc6ffc93a205bd66c-a8f8f5138ef1e3b356bd5078a32694625e8846a5" | sha1sum -c - \
  && curl -fLo $HOME/.bosh/downloads/9f38ce54a268f99e2d6e37723d813fce3457b970-cbbf73c102b1f27d3db15d95bc971b2b4995c78e https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-google-cpi-release?v=27.0.0 \
    && echo "cbbf73c102b1f27d3db15d95bc971b2b4995c78e  $HOME/.bosh/downloads/9f38ce54a268f99e2d6e37723d813fce3457b970-cbbf73c102b1f27d3db15d95bc971b2b4995c78e" | sha1sum -c -
ADD boshua /bin/boshua
