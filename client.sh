#!/bin/bash

# if ever needed
function help() {
  echo "bash cli clent for jackdwyer.org"
  echo "./client.sh <host>"
}

HOST=${1}

if [[ -z ${HOST} ]]; then
  HOST=http://localhost:5000/upload
fi

# Melbourne
LATITUDE=-37.814
LONGITUDE=144.96332

# LATITUDE=40.70574
# LONGITUDE=-73.94249
IMAGE=test_data/test1.jpg
curl -v -i -X POST \
-H "Content-Type: multipart/form-data"  \
-F "img=@${IMAGE}" \
-F "latitude=${LATITUDE}" \
-F "longitude=${LONGITUDE}" \
${HOST}

# LATITUDE=40.70574
# LONGITUDE=-73.94249
IMAGE=test_data/test2.jpg
curl -v -i -X POST \
-H "Content-Type: multipart/form-data"  \
-F "img=@${IMAGE}" \
-F "latitude=${LATITUDE}" \
-F "longitude=${LONGITUDE}" \
${HOST}

LATITUDE=40.70574
LONGITUDE=-73.94249
IMAGE=test_data/test3.jpg
curl -v -i -X POST \
-H "Content-Type: multipart/form-data"  \
-F "img=@${IMAGE}" \
-F "latitude=${LATITUDE}" \
-F "longitude=${LONGITUDE}" \
${HOST}
