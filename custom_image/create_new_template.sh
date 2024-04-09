#!/bin/bash

# Path to the existing template JSON configuration
#CONFIG_JSON="existing_template_config.json"

# Define new custom OS image (this should be set or passed into the script)
NEW_IMAGE=$PACKER_IMAGE_NAME
# GCP_PROJECT_ID=${{ vars.GCP_PROJECT_ID }}
INSTANCE_GROUP_NAME=$INSTANCE_GROUP_NAME

# Extracting information from the JSON file
MACHINE_TYPE=$REGION
DISK_SIZE_GB=$DISK_SIZE_GB
DISK_TYPE=$DISK_TYPE
#DISK_INTF_TYPE=$(jq -r '.properties.disks[0].interface' "$CONFIG_JSON")

#NETWORK=$(jq -r '.properties.networkInterfaces[0].network' "$CONFIG_JSON" | awk -F'/' '{print $NF}')
#SUBNET=$(jq -r '.properties.networkInterfaces[0].subnetwork' "$CONFIG_JSON" | awk -F'/' '{print $NF}')
REGION=$REGION
NETWORK=$NETWORK
SUBNET=$SUBNET
# Assuming tags are listed under properties.tags.items in the JSON
TAGS=$TAGS

echo "$REGION"
echo "$NETWORK"
echo "$SUBNET"
echo "$TAGS"g
#echo $DISK_INTF_TYPE


# Extract the startup script from the JSON and save it to a file
# STARTUP_SCRIPT_VALUE=$(jq -r '.properties.metadata.items[] | select(.key=="startup-script") | .value' "$CONFIG_JSON")
# echo $STARTUP_SCRIPT_VALUE
# echo "$STARTUP_SCRIPT_VALUE" > startup-script.sh

cat << 'EOF' > startup-script.sh
#!/bin/bash

get_secret() {
  gcloud secrets versions access latest --secret="$1"
}
if [ ! -f /etc/webapp.flag ]; then
  DB_USER=$(get_secret "db_user")
  DB_PASSWORD=$(get_secret "db_password")
  DB_NAME=$(get_secret "db_name")
  DB_HOST=$(get_secret "db_host")

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


# Service account and scopes
SERVICE_ACCOUNT_EMAIL=$SERVICE_ACCOUNT_EMAIL
SCOPES=$SCOPES
NEW_INSTANCE_TEMPALTE="webapp-template-$(date +%Y%m%d%H%M%S)"
# Constructing the gcloud command
GCLOUD_CMD="gcloud compute instance-templates create $NEW_INSTANCE_TEMPALTE \
    --machine-type=$MACHINE_TYPE \
    --image=$NEW_IMAGE \
    --image-project=$GCP_PROJECT_ID \
    --boot-disk-size=${DISK_SIZE_GB}GB \
    --boot-disk-type=$DISK_TYPE \
    --network-interface=network=$NETWORK,subnet=$SUBNET \
    --instance-template-region=$REGION \
    --tags=$TAGS \
    --address="" \
    --metadata-from-file=startup-script=startup-script.sh \
    --service-account=$SERVICE_ACCOUNT_EMAIL \
    --scopes=$SCOPES"

echo "$GCLOUD_CMD"
# Execute the constructed gcloud command
eval $GCLOUD_CMD

echo "New template created: $NEW_INSTANCE_TEMPALTE"

# Construct the new instance template URL
NEW_INSTANCE_TEMPLATE_URL="https://www.googleapis.com/compute/v1/projects/$GCP_PROJECT_ID/regions/$REGION/instanceTemplates/$NEW_INSTANCE_TEMPALTE"
echo "NEW_INSTANCE_TEMPLATE_URL $NEW_INSTANCE_TEMPLATE_URL"

TEMPLATE_READY="false"
while [[ $TEMPLATE_READY == "false" ]]; do
  output=$(gcloud compute instance-templates describe $NEW_INSTANCE_TEMPLATE_URL --format="get(name)" 2>&1)
    if [[ $output != *"ERROR"* ]]; then
    echo "Instance template is ready: $output"
    TEMPLATE_READY="true"
  else
    echo "Waiting for instance template to be ready..."
    sleep 10
  fi
done

# Update the managed instance group to use the new template
gcloud compute instance-groups managed set-instance-template $INSTANCE_GROUP_NAME \
  --template=$NEW_INSTANCE_TEMPLATE_URL \
  --region=$REGION

# #Update the managed instance group to use the new template
# gcloud compute instance-groups managed set-instance-template webapp-group \
#   --template="$NEW_INSTANCE_TEMPLATE_URL" \
#   --region=us-east1

echo "Updated instance group ${INSTANCE_GROUP_NAME} to use new template: ${NEW_INSTANCE_TEMPLATE_URL}"

gcloud compute instance-groups managed rolling-action start-update $INSTANCE_GROUP_NAME \
  --version=template=$NEW_INSTANCE_TEMPLATE_URL \
  --region=$REGION \
  --max-unavailable=0

# gcloud compute instance-groups managed rolling-action start-update webapp-group \
#   --version=template="$NEW_INSTANCE_TEMPLATE_URL" \
#   --region=us-east1 \
#   --max-unavailable=0
