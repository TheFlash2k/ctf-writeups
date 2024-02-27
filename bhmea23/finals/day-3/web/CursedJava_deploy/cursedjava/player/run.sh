#!/bin/bash
# * Before you read this source code, please note that:
# * I would like to apologize for this crime against humanity
# * This is what running a Java app in nsjail looks like

export JAVA_HOME=/usr/local/openjdk-17
export JAVA_VERSION=17.0.2
export PATH=$JAVA_HOME/bin:$PATH

redis-server /etc/redis/redis.conf >/dev/null 2>&1 &
java -jar /app/app.jar >/app/app.log 2>&1 &
nginx -e /dev/null &
while ! nc -z 127.0.0.1 8080; do
    sleep 0.1
done
socat - TCP:127.0.0.1:9999,forever
