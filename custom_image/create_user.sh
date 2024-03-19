echo Creating group csye6225
sudo groupadd csye6225
echo Creating user csye6225 with no login shell
sudo useradd -r -g csye6225 -s /usr/sbin/nologin csye6225
echo User and group csye6225 created successfully
sudo mkdir -p /var/log/webapp
sudo chown csye6225:csye6225 /var/log/webapp
sudo chmod 755 /var/log/webapp
echo /var/log/webapp created successfully
