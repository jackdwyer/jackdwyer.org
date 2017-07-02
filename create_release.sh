#!/bin/bash
set -euo pipefail

TAG=${1-}
TAG_OBJECT=${2-}

if [[ -z ${TAG} ]] || [[ -z ${TAG_OBJECT} ]]; then
  echo "./create release <tag> <sha>"
  exit 1
fi

TAG_MESSAGE="CI GENERATED"
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
RELEASE_NAME="${TAG_OBJECT}"
RELEASE_DESC="${RELEASE_NAME}"

tag_body=$(cat  << EOF
{
  "tag": "${TAG}",
  "message": "${TAG_MESSAGE}",
  "object": "${TAG_OBJECT}",
  "type": "commit",
  "tagger": {
    "name": "CI",
    "email": "jackjack.dwyer@gmail.com",
    "date": "${DATE}"
  }
}
EOF
)

curl -X POST https://api.github.com/repos/jackdwyer/jackdwyer.org/git/tags \
-u "jackdwyer:${GITHUB_API_TOKEN}" \
-d "${tag_body}"

release_body=$(cat  << EOF
{
  "tag_name": "${TAG}",
  "target_commitsh": "master",
  "name": "${RELEASE_NAME}",
  "body": "${RELEASE_DESC}",
  "draft": false,
  "prerelease": false
}
EOF
)

asset_upload_url=$(curl -s -X POST https://api.github.com/repos/jackdwyer/jackdwyer.org/releases \
-u "jackdwyer:${GITHUB_API_TOKEN}" \
-d "${release_body}" | jq -r -M .upload_url | cut -d"{" -f1)

if ! [[ -f bin/jackdwyer ]]; then
  make release
fi

curl -X POST "${asset_upload_url}?name=jackdwyer" \
-u "jackdwyer:${GITHUB_API_TOKEN}" \
-H "Content-Type: application/octet-stream" \
--data-binary @bin/jackdwyer
