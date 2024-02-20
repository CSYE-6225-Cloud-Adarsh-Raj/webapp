#!/bin/bash

SERVICE_FILE="/etc/systemd/system/webapp.service"

# Replace placeholders with actual values
sed -i "s|\\${DB_USER}|$DB_USER|g" $SERVICE_FILE
sed -i "s|\\${DB_HOST}|$DB_HOST|g" $SERVICE_FILE
sed -i "s|\\${DB_PASSWORD}|$DB_PASSWORD|g" $SERVICE_FILE
sed -i "s|\\${DB_NAME}|$DB_NAME|g" $SERVICE_FILE
