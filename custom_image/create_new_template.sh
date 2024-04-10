#!/bin/bash
set -e

get_secret() {
  gcloud secrets versions access latest --secret="$1"
}
DB_USER=$(get_secret "db-user")
DB_PASSWORD=$(get_secret "db-password")
DB_NAME=$(get_secret "db-name")
DB_HOST=$(get_secret "db-host")
BOOT_DISK_KMS_KEY=$(get_secret "vm-key")

cat << 'EOF' > startup-script.sh
#!/bin/bash
if [ ! -f /etc/webapp.flag ]; then
  echo "DB_USER=$DB_USER" > /etc/webapp.env 
  echo "DB_PASSWORD=$DB_PASSWORD" >> /etc/webapp.env
  echo "DB_NAME=$DB_NAME" >> /etc/webapp.env 
  echo "DB_HOST=$DB_HOST" >> /etc/webapp.env 
  sudo chown csye6225:csye6225 /etc/systemd/system/webapp.service 
  sudo touch /etc/webapp.flag 
else 
  echo "/etc/webapp.flag exists, skipping script execution." 
fi
EOF

NEW_INSTANCE_TEMPALTE="webapp-template-$(date +%Y%m%d%H%M%S)"
# Constructing the gcloud command
GCLOUD_CMD="gcloud compute instance-templates create $NEW_INSTANCE_TEMPALTE \
    --region=$REGION \
    --machine-type=$MACHINE_TYPE \
    --image=$PACKER_IMAGE_NAME \
    --image-project=$GCP_PROJECT_ID \
    --boot-disk-size=${DISK_SIZE_GB}GB \
    --boot-disk-type=$DISK_TYPE \
    --network-interface=network=$NETWORK,subnet=$SUBNET \
    --instance-template-region=$REGION \
    --tags=$TAGS \
    --address="" \
    --metadata-from-file=startup-script=startup-script.sh \
    --service-account=$SERVICE_ACCOUNT_EMAIL \
    --project=$GCP_PROJECT_ID \
    --boot-disk-kms-key= $BOOT_DISK_KMS_KEY \
    --scopes=$SCOPES"

echo "$GCLOUD_CMD"
# Execute the constructed gcloud command
eval $GCLOUD_CMD

echo "New template created: $NEW_INSTANCE_TEMPALTE"

NEW_INSTANCE_TEMPLATE_URL=$(gcloud compute instance-templates describe $NEW_INSTANCE_TEMPALTE \
--region=$REGION \
--format="value(self_link)")
  
# Update the managed instance group to use the new template
gcloud compute instance-groups managed set-instance-template $INSTANCE_GROUP_NAME \
  --template=$NEW_INSTANCE_TEMPLATE_URL \
  --region=$REGION

echo "Updated instance group ${INSTANCE_GROUP_NAME} to use new template: ${NEW_INSTANCE_TEMPLATE_URL}"

gcloud compute instance-groups managed rolling-action start-update $INSTANCE_GROUP_NAME \
  --version=template=$NEW_INSTANCE_TEMPLATE_URL \
  --region=$REGION \
  --max-unavailable=0
