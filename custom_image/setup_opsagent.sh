echo Checking if /etc/google-cloud-ops-agent/conf.d exists...
if [ ! -d /etc/google-cloud-ops-agent/conf.d ]; then
  echo Directory /etc/google-cloud-ops-agent/conf.d does not exist. Creating it...
  sudo mkdir -p /etc/google-cloud-ops-agent/conf.d
else
  echo Directory /etc/google-cloud-ops-agent/conf.d exists.
fi

echo Moving /tmp/config.yaml to /etc/google-cloud-ops-agent/config.yaml
sudo mv /tmp/config.yaml /etc/google-cloud-ops-agent/config.yaml
echo Move operation completed.