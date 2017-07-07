#!/bin/bash

USER=dwyerj
HOST=dev.jackdwyer.org

TIMESTAMP=$1
if [[ -z ${TIMESTAMP} ]]; then
  echo "./deploy.sh <timestamp>"
  exit 1
fi

scp -i yar_my_key jackdwyer.org.service templates/* ${USER}@${HOST}:/opt/jackdwyer/templates/

ssh -i yar_my_key ${USER}@${HOST} "
  sudo systemctl stop jackdwyer.org
  curl -L -o /opt/jackdwyer/jackdwyer 'https://github.com/jackdwyer/jackdwyer.org/releases/download/${TIMESTAMP}/jackdwyer'
  chmod +x /opt/jackdwyer/jackdwyer
  sudo systemctl daemon-reload
  sudo systemctl start jackdwyer.org
"
