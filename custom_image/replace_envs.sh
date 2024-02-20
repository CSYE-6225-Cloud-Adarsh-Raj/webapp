#!/bin/bash

TMP_SERVICE_FILE="/tmp/webapp.service"

# Make a copy to /tmp where we have write access
sudo cp /etc/systemd/system/webapp.service $TMP_SERVICE_FILE

# Replace placeholders with actual values in the temporary file
sed -i "s|\\${DB_USER}|$DB_USER|g" $TMP_SERVICE_FILE
sed -i "s|\\${DB_HOST}|$DB_HOST|g" $TMP_SERVICE_FILE
sed -i "s|\\${DB_PASSWORD}|$DB_PASSWORD|g" $TMP_SERVICE_FILE
sed -i "s|\\${DB_NAME}|$DB_NAME|g" $TMP_SERVICE_FILE

sudo mv $TMP_SERVICE_FILE /etc/systemd/system/webapp.service