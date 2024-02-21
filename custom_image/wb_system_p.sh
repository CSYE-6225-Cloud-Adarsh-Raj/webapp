echo Moving /tmp/webapp to /usr/local/bin
sudo mv /tmp/webapp /usr/local/bin/webapp
sudo chmod +x /usr/local/bin/webapp
sudo chown csye6225:csye6225 /usr/local/bin/webapp
echo Setting SELinux context for /usr/local/bin/webapp
sudo semanage fcontext -a -t bin_t '/usr/local/bin/webapp'
sudo restorecon -v '/usr/local/bin/webapp'
echo Listing contents of /usr/local/bin
ls -la /usr/local/bin/
