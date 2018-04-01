#!/bin/bash

set -eu -o pipefail

git clone --quiet "file://$PWD/index" index-out

mkdir -p index-out/$storage

tarball_path=$( echo $PWD/compiled-release/*.tgz )
tarball_name="$( basename "$tarball_path" )"
metalink_path="index-out/$storage/compiled-release.meta4"

meta4 create --metalink="$metalink_path"
meta4 set-published --metalink="$metalink_path" "$( date -u +%Y-%m-%dT%H:%M:%SZ )"
meta4 import-file --metalink="$metalink_path" --file="$tarball_name" --version="$release_version" "$tarball_path"

sha256=$( meta4 file-hash --metalink="$metalink_path" sha-256 )
path1=$( echo "$sha256" | cut -c-2 )
path2=$( echo "$sha256" | cut -c3-4 )

export AWS_ACCESS_KEY_ID="$s3_access_key"
export AWS_SECRET_ACCESS_KEY="$s3_secret_key"

meta4 file-upload --metalink="$metalink_path" --file="$tarball_name" "$tarball_path" "s3://$s3_host/$s3_bucket/$path1/$path2/$sha256"

echo "$context" > "index-out/$storage/compiled-release.json"

mv compiled-release/compilation.json "index-out/$storage/compilation.json"

cd index-out

if [[ -z "$( git status --porcelain )" ]]; then
  exit
fi

export GIT_COMMITTER_EMAIL="concourse.ci@localhost"
export GIT_COMMITTER_NAME="Concourse"

git config --global user.email "$GIT_COMMITTER_EMAIL"
git config --global user.name "$GIT_COMMITTER_NAME"

git add .

git commit -m "$tarball_name"