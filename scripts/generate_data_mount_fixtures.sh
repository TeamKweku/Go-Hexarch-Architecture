#!/bin/bash

set -eou pipefail

echo "Generating development data mount fixtures..."

echo "Checking for existence of host data directory '${CODE_ODESSEY_HOST_DATA_DIR}'..."
if [ -d "${CODE_ODESSEY_HOST_DATA_DIR}" ]; then
  echo "Host data directory already exists. Skipping generation."
  echo
  exit 0
fi

echo "Creating gitignored host data directory..."
mkdir "${CODE_ODESSEY_HOST_DATA_DIR}"
echo "Host data directory created."

echo "Generating development-only JWT RSA private key..."
ssh-keygen -t rsa -b 4096 -m PEM -f "${CODE_ODESSEY_HOST_DATA_DIR}/${CODE_ODESSEY_JWT_RSA_PRIVATE_KEY_PEM_BASENAME}" -N ""
echo "Development JWT key generated."

echo "Data mount fixtures generated successfully at '${CODE_ODESSEY_HOST_DATA_DIR}.'"
