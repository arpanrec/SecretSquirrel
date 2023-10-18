#!/usr/bin/env bash
set -e

SECURE_SERVER_DIR=/opt/secureserver
SECURE_SERVER_USER=secureserver
SECURE_SERVER_GROUP=secureserver
SECURE_SERVER_SYSTEMD_SERVICE_NAME=secureserver.service
sudo mkdir -p "${SECURE_SERVER_DIR}"

echo "Installing Secure Server"
echo "This script requires sudo privileges"
echo "This script will install the Secure Server as a systemd service"
echo "This script will install the Secure Server to /opt/secureserver"
echo "This script will create a user and group named secureserver"
echo "This script will create a systemd service named secureserver.service"
echo "Sourcing .env file if it exists"

if [ -f .env ]; then
    # shellcheck disable=SC1091
    source .env
    echo "Writing INIT_* variables to file"
    echo -n "${INIT_GPG_PRIVATE_KEY_PASSWORD}" | sudo tee "${SECURE_SERVER_DIR}/INIT_GPG_PRIVATE_KEY_PASSWORD" >/dev/null
    echo -n "${INIT_OPENSSL_ROOT_CA_KEY_PASSWORD}" | sudo tee "${SECURE_SERVER_DIR}/INIT_OPENSSL_ROOT_CA_KEY_PASSWORD" >/dev/null
    echo -n "${INIT_OPENSSL_ROOT_CA_CERT_CONTENT}" | sudo tee "${SECURE_SERVER_DIR}/INIT_OPENSSL_ROOT_CA_CERT" >/dev/null
    echo -n "${INIT_OPENSSL_ROOT_CA_KEY_CONTENT}" | sudo tee "${SECURE_SERVER_DIR}/INIT_OPENSSL_ROOT_CA_KEY" >/dev/null
    echo -n "${INIT_GPG_PRIVATE_KEY_CONTENT}" | sudo tee "${SECURE_SERVER_DIR}/INIT_GPG_PRIVATE_KEY" >/dev/null
    echo -n "${INIT_GPG_PRIVATE_KEY_CONTENT}" | sudo tee "${SECURE_SERVER_DIR}/INIT_GPG_PUBLIC_KEY" >/dev/null
    echo -n "${INIT_GPG_PRIVATE_KEY_FINGERPRINT}" | sudo tee "${SECURE_SERVER_DIR}/INIT_GPG_PRIVATE_KEY_FINGERPRINT" >/dev/null
    echo -n "${INIT_EMAIL_ID}" | sudo tee "${SECURE_SERVER_DIR}/INIT_EMAIL_ID" >/dev/null
else
    echo "No .env file found"
    exit 1
fi

sudo systemctl disable --now "${SECURE_SERVER_SYSTEMD_SERVICE_NAME}" || true

sudo rm -rf "${SECURE_SERVER_DIR}/secureserver" ./secureserver "${SECURE_SERVER_DIR}/config.json"
go build -o "secureserver" ./main.go
sudo mv ./secureserver "${SECURE_SERVER_DIR}/secureserver"
sudo cp ./config-prod.json "${SECURE_SERVER_DIR}/config.json"
sudo setcap 'cap_net_bind_service=+ep' "${SECURE_SERVER_DIR}/secureserver"
sudo userdel -r "${SECURE_SERVER_USER}" || true
sudo groupdel "${SECURE_SERVER_GROUP}" || true
sudo groupadd --system "${SECURE_SERVER_GROUP}"
sudo useradd -s /bin/false --home-dir "/home/${SECURE_SERVER_USER}" --no-create-home \
    --system --gid "${SECURE_SERVER_GROUP}" "${SECURE_SERVER_USER}" || true
sudo chown -R "${SECURE_SERVER_USER}":"${SECURE_SERVER_GROUP}" "${SECURE_SERVER_DIR}"

sudo ufw allow 80/tcp
sudo ufw allow 433/tcp
sudo systemctl enable --now ufw
sudo sudo systemctl restart ufw
sudo ufw --force enable
sudo ufw reload

sudo docker run --rm \
    -v /etc/letsencrypt:/etc/letsencrypt \
    -v /var/log/letsencrypt:/var/log/letsencrypt \
    -p 80:80 \
    certbot/certbot \
    certonly \
    --standalone \
    --non-interactive \
    --agree-tos \
    --email "${INIT_EMAIL_ID}" \
    --domains 172-105-49-235.ip.linodeusercontent.com \
    --preferred-challenges http-01 >/dev/null

cat <<EOF | sudo tee /etc/systemd/system/"${SECURE_SERVER_SYSTEMD_SERVICE_NAME}" >/dev/null
[Unit]
Description=Secure Server
After=network.target

[Service]
Type=simple
User=${SECURE_SERVER_USER}
Group=${SECURE_SERVER_GROUP}
WorkingDirectory=${SECURE_SERVER_DIR}
ExecStart=${SECURE_SERVER_DIR}/secureserver
#Environment=SECURE_SERVER_CONFIG_FILE_PATH=${SECURE_SERVER_DIR}/config.json
Restart=always

[Install]
WantedBy=multi-user.target
EOF

echo "Starting Secure Server"
sudo systemctl daemon-reload
sudo systemctl enable --now "${SECURE_SERVER_SYSTEMD_SERVICE_NAME}"
sudo systemctl status "${SECURE_SERVER_SYSTEMD_SERVICE_NAME}"
