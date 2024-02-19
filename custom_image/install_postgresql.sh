#!/bin/bash
DB_USER=${TEST_DB_USER}
DB_PASSWORD=${TEST_DB_PASSWORD}
DB_NAME=${TEST_DB_NAME}

sudo dnf -y install https://download.postgresql.org/pub/repos/yum/reporpms/EL-8-x86_64/pgdg-redhat-repo-latest.noarch.rpm
sudo dnf -qy module disable postgresql
sudo dnf -y install postgresql16-server
sudo /usr/pgsql-16/bin/postgresql-16-setup initdb
sudo systemctl enable postgresql-16
sudo systemctl start postgresql-16
sudo -u postgres psql -c "SELECT version();"

# Create a PostgreSQL user and database, then grant privileges
sudo -u postgres psql -c "CREATE USER ${DB_USER} WITH ENCRYPTED PASSWORD '${DB_PASSWORD}';"
sudo -u postgres psql -c "CREATE DATABASE ${DB_NAME};"
sudo -u postgres psql -c "ALTER ROLE ${DB_USER} SUPERUSER;"