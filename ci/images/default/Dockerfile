FROM golang:1.10-stretch
RUN true \
  && apt-get update \
  && apt-get install -y curl git openssl build-essential zlibc zlib1g-dev ruby ruby-dev libxslt1-dev libxml2-dev libssl-dev readline-common libreadline-dev libyaml-dev libsqlite3-dev sqlite3 \
  && curl -fLo /usr/local/bin/bosh https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-3.0.1-linux-amd64 \
    && echo "ccc893bab8b219e9e4a628ed044ebca6c6de9ca0  /usr/local/bin/bosh" | sha1sum -c - \
    && chmod +x /usr/local/bin/bosh \
  && curl -fLo /usr/local/bin/meta4 https://github.com/dpb587/metalink/releases/download/v0.2.0/meta4-0.2.0-linux-amd64 \
    && echo "9b7afe863d95ed8d99886dbe1878619a71b97725  /usr/local/bin/meta4" | sha1sum -c - \
    && chmod +x /usr/local/bin/meta4 \
  && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
