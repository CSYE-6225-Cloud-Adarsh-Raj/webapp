sudo mv /tmp/restart_webapp.sh /usr/local/bin/restart_webapp.sh
sudo chmod +x /usr/local/bin/restart_webapp.sh
sudo chown csye6225:csye6225 /usr/local/bin/restart_webapp.sh
sudo semanage fcontext -a -t bin_t '/usr/local/bin/restart_webapp.sh'
sudo restorecon -v '/usr/local/bin/restart_webapp.sh'
