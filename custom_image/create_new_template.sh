#!/bin/bash

# Path to the existing template JSON configuration
CONFIG_JSON="existing_template_config.json"

# Define new custom OS image (this should be set or passed into the script)
NEW_IMAGE="${PACKER_IMAGE_NAME}"
NEW_IMAGE_PROJECT_ID="csye6225-dev-414220"

# Extracting information from the JSON file
MACHINE_TYPE=$(jq -r '.properties.machineType' "$CONFIG_JSON")
DISK_SIZE_GB=$(jq -r '.properties.disks[0].initializeParams.diskSizeGb' "$CONFIG_JSON")
DISK_TYPE=$(jq -r '.properties.disks[0].initializeParams.diskType' "$CONFIG_JSON")
#DISK_INTF_TYPE=$(jq -r '.properties.disks[0].interface' "$CONFIG_JSON")

#NETWORK=$(jq -r '.properties.networkInterfaces[0].network' "$CONFIG_JSON" | awk -F'/' '{print $NF}')
#SUBNET=$(jq -r '.properties.networkInterfaces[0].subnetwork' "$CONFIG_JSON" | awk -F'/' '{print $NF}')
REGION=$(jq -r '.region' "$CONFIG_JSON" | awk -F'/' '{print $NF}')
NETWORK=$(jq -r '.properties.networkInterfaces[0].network' "$CONFIG_JSON")
SUBNET=$(jq -r '.properties.networkInterfaces[0].subnetwork' "$CONFIG_JSON")
# Assuming tags are listed under properties.tags.items in the JSON
TAGS=$(jq -r '.properties.tags.items | join(",")' "$CONFIG_JSON")

echo $REGION
echo $NETWORK
echo $SUBNET
echo $TAGS
#echo $DISK_INTF_TYPE


# Extract the startup script from the JSON and save it to a file
STARTUP_SCRIPT_VALUE=$(jq -r '.properties.metadata.items[] | select(.key=="startup-script") | .value' "$CONFIG_JSON")
echo "$STARTUP_SCRIPT_VALUE" > startup-script.sh

# Service account and scopes
SERVICE_ACCOUNT_EMAIL=$(jq -r '.properties.serviceAccounts[0].email' "$CONFIG_JSON")
SCOPES=$(jq -r '.properties.serviceAccounts[0].scopes | join(",")' "$CONFIG_JSON")
NEW_TEMPLATE_NAME="webapp-template-$(date +%Y%m%d%H%M%S)"
# Constructing the gcloud command
GCLOUD_CMD="gcloud compute instance-templates create $NEW_TEMPLATE_NAME \
    --machine-type=$MACHINE_TYPE \
    --image=$NEW_IMAGE \
    --image-project=$NEW_IMAGE_PROJECT_ID \
    --boot-disk-size=${DISK_SIZE_GB}GB \
    --boot-disk-type=$DISK_TYPE \
    --network-interface=network=$NETWORK,subnet=$SUBNET \
    --instance-template-region=$REGION \
    --tags=$TAGS \
    --address="" \
    --metadata-from-file=startup-script=startup-script.sh \
    --service-account=$SERVICE_ACCOUNT_EMAIL \
    --scopes=$SCOPES"


# Echo the command for testing and verification
echo "$GCLOUD_CMD"

# Execute the constructed gcloud command
eval "$GCLOUD_CMD"

echo "New template created: $NEW_TEMPLATE_NAME"


# REGION="us-east1"
# INSTANCE_GROUP_NAME="webapp-group" # Make sure to replace this with your actual instance group name
# INSTANCE_GROUP_NAME: ${{ secrets.DEFAULT_INSTANCE_GROUP_NAME }}
# The rest of your script up to the echo "New template created: $NEW_TEMPLATE_NAME"

# Construct the new instance template URL
NEW_INSTANCE_TEMPLATE_URL="https://www.googleapis.com/compute/v1/projects/${NEW_IMAGE_PROJECT_ID}/regions/${REGION}/instanceTemplates/${NEW_TEMPLATE_NAME}"
echo "NEW_INSTANCE_TEMPLATE_URL ${NEW_INSTANCE_TEMPLATE_URL} 

# Update the managed instance group to use the new template
gcloud compute instance-groups managed set-instance-template "${INSTANCE_GROUP_NAME}" \
  --template="${NEW_INSTANCE_TEMPLATE_URL}" \
  --region="${REGION}"

echo "Updated instance group ${INSTANCE_GROUP_NAME} to use new template: ${NEW_INSTANCE_TEMPLATE_URL}"
