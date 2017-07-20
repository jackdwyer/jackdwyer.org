#!/bin/bash

LATITUDE=40.70574
LONGITUDE=-73.94249
IMAGE=test_data/test1.jpg
curl -v -i -X POST \
-H "Content-Type: multipart/form-data"  \
-F "img=@${IMAGE}" \
-F "latitude=${LATITUDE}" \
-F "longitude=${LONGITUDE}" \
http://localhost:5000/upload

LATITUDE=40.70574
LONGITUDE=-73.94249
IMAGE=test_data/test2.jpg
curl -v -i -X POST \
-H "Content-Type: multipart/form-data"  \
-F "img=@${IMAGE}" \
-F "latitude=${LATITUDE}" \
-F "longitude=${LONGITUDE}" \
http://localhost:5000/upload

LATITUDE=40.70574
LONGITUDE=-73.94249
IMAGE=test_data/test3.jpg
curl -v -i -X POST \
-H "Content-Type: multipart/form-data"  \
-F "img=@${IMAGE}" \
-F "latitude=${LATITUDE}" \
-F "longitude=${LONGITUDE}" \
http://localhost:5000/upload
