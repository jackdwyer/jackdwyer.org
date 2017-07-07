#!/bin/bash

USER=dwyerj
HOST=dev.jackdwyer.org

TIMESTAMP=$1
if [[ -z ${TIMESTAMP} ]]; then
  echo "./deploy.sh <timestamp>"
  exit 1
fi

SSH_ARGS="-i yar_my_key -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"

scp ${SSH_ARGS} jackdwyer.org.service templates/* ${USER}@${HOST}:/opt/jackdwyer/templates/

ssh ${SSH_ARGS} ${USER}@${HOST} "
  sudo systemctl stop jackdwyer.org
  curl -L -o /opt/jackdwyer/jackdwyer 'https://github.com/jackdwyer/jackdwyer.org/releases/download/${TIMESTAMP}/jackdwyer'
  chmod +x /opt/jackdwyer/jackdwyer
  sudo systemctl daemon-reload
  sudo systemctl start jackdwyer.org
"
