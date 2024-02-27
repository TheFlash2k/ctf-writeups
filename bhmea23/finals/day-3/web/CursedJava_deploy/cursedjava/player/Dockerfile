# Java app build
FROM maven:3.8.4-openjdk-17-slim AS build
WORKDIR /app
COPY ./api/pom.xml .
RUN mvn dependency:go-offline
COPY ./api/src ./src
RUN mvn package -DskipTests

# Building front
FROM node:18-slim AS frontbuild
WORKDIR /usr/src/app
COPY ./front/package*.json ./
RUN npm install
COPY ./front/ .
RUN npm run build

# Base image for container, contains everything
FROM openjdk:17-jdk-slim AS final

WORKDIR /app
ARG JAR_FILE=/app/target/*.jar
COPY --from=build ${JAR_FILE} app.jar

# Installing nginx
RUN apt update && \
    apt install -y curl gnupg2 ca-certificates lsb-release debian-archive-keyring wget && \
    curl https://nginx.org/keys/nginx_signing.key | gpg --dearmor | tee /usr/share/keyrings/nginx-archive-keyring.gpg >/dev/null && \
    echo "deb [signed-by=/usr/share/keyrings/nginx-archive-keyring.gpg] http://nginx.org/packages/debian `lsb_release -cs` nginx" | tee /etc/apt/sources.list.d/nginx.list && \
    apt update && \
    apt install -y nginx socat sqlite3 redis-server netcat-openbsd && \
    rm -rf /var/lib/apt/lists/*

# Serving front
COPY --from=frontbuild /usr/src/app/dist /usr/share/nginx/html
COPY ./front/nginx.conf /etc/nginx/nginx.conf
COPY ./redis.conf /etc/redis/redis.conf

FROM pwn.red/jail:latest

COPY --from=final / /srv
COPY hook.sh /jail/hook.sh
COPY ./run.sh /srv/app/run
RUN chmod +x /srv/app/run /jail/hook.sh && \
    touch /srv/app/app.log && \
    mkdir /srv/app/data && \
    chmod 777 -R /srv/app/ && \
    chmod 777 /srv/etc && \
    chmod 777 -R /srv/etc/redis && \
    echo "127.0.0.1 localhost" >> /srv/etc/hosts
ENV JAIL_PIDS=0 JAIL_CPU=0 JAIL_MEM=0 JAIL_TIME=0 JAIL_TMP_SIZE=100M
