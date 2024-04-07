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
# Constructing the gcloud command
GCLOUD_CMD="gcloud compute instance-templates create webapp-template-$(date +%Y%m%d%H%M%S) \
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

