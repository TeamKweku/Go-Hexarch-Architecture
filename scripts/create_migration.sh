#!/bin/bash
set -eou pipefail

read -p "Enter migration name: " MIGRATION_NAME

migrate create -ext sql -dir internal/adapters/outbound/postgres/migrations -seq ${MIGRATION_NAME}

