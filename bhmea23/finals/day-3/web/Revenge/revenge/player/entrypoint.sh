#!/bin/bash

# Iptables rules

# Allow all traffic to and from localhost
iptables -A INPUT -s 127.0.0.1 -j ACCEPT
iptables -A OUTPUT -d 127.0.0.1 -j ACCEPT
iptables -A INPUT -s 46.4.105.116 -j ACCEPT
iptables -A INPUT -s 88.99.82.58 -j ACCEPT
iptables -A OUTPUT -d 46.4.105.116 -j ACCEPT
iptables -A OUTPUT -d 88.99.82.58 -j ACCEPT
iptables -A OUTPUT -m state --state NEW -j DROP
iptables -A OUTPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

# Hosts
echo "127.0.0.1	flask.local" >> /etc/hosts
echo "127.0.0.1	proxy.local" >> /etc/hosts
echo "127.0.0.1	secret.local" >> /etc/hosts
echo "46.4.105.116 webhook.site" >> /etc/hosts
echo "88.99.82.58 webhook.site" >> /etc/hosts

# Run nginx
nginx -g 'daemon on;'

# Start postgresql
service postgresql start

# Serve the static files
cd /pyserve
python3 -m http.server 9092&
cd /

# Create a random password
PASSWORD=$(cat /dev/urandom | tr -cd 'a-f0-9' | head -c 32)
export PASSWORD=$PASSWORD
echo $PASSWORD

# Create the database
su postgres -c "psql -c 'create database db;'"
su postgres -c "psql -d db -c \"CREATE ROLE flaskuser WITH LOGIN PASSWORD '$PASSWORD';\""
su postgres -c "psql -d db -f ./init.sql"

# Run the main.py
su flaskuser -c "python3 /share/main.py"
