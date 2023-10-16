#!/usr/bin/env bash
set -ex
SECURE_SERVER_DIR=/opt/secureserver
SECURE_SERVER_USER=secureserver
SECURE_SERVER_GROUP=secureserver
SECURE_SERVER_SYSTEMD_SERVICE_NAME=secureserver.service
sudo systemctl disable --now "${SECURE_SERVER_SYSTEMD_SERVICE_NAME}" || true

sudo mkdir -p "${SECURE_SERVER_DIR}"
go build ./main.go
sudo mv main "${SECURE_SERVER_DIR}/secureserver"
sudo rm -rf "${SECURE_SERVER_DIR}/config.json"
sudo cp ./config.json "${SECURE_SERVER_DIR}/config.json"
sudo userdel -r "${SECURE_SERVER_USER}" || true
sudo groupdel "${SECURE_SERVER_GROUP}" || true
sudo groupadd --system "${SECURE_SERVER_GROUP}"
sudo useradd -s /bin/false --home-dir "/home/${SECURE_SERVER_USER}" --no-create-home "${SECURE_SERVER_USER}" \
    --system --gid "${SECURE_SERVER_GROUP}" "${SECURE_SERVER_USER}" || true
sudo chown -R "${SECURE_SERVER_USER}":"${SECURE_SERVER_GROUP}" "${SECURE_SERVER_DIR}/secureserver"
