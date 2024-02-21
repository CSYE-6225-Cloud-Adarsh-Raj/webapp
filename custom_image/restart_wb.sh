sudo mv /tmp/restart_webapp.service /etc/systemd/system/restart_webapp.service
sudo chown csye6225:csye6225 /etc/systemd/system/restart_webapp.service
sudo systemctl enable restart_webapp.service
