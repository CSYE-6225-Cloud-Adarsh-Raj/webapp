#!/bin/bash
curl -LO https://golang.org/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh
#echo 'export PATH=$PATH:/usr/local/go/bin' >> /home/packer/.bash_profile

