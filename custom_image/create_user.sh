echo Creating group csye6225
sudo groupadd csye6225
echo Creating user csye6225 with no login shell
sudo useradd -r -g csye6225 -s /usr/sbin/nologin csye6225
echo User and group csye6225 created successfully
