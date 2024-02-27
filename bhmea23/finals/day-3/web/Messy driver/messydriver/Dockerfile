FROM golang:1.20.6 AS builder

RUN apt-get update -y && apt install sqlite3 -y && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./main.go /app/main.go
COPY ./idp /app/idp/
COPY ./fileserver /app/fileserver/
COPY ./proxy /app/proxy/
COPY ./utils /app/utils/
COPY ./idp /app/idp
COPY ./setup.sh /app/setup.sh
RUN go build -o server ./main.go
# ! Most likely reduntant due to the jail's run.sh!!! pay attention!
RUN bash /app/setup.sh

# ----------------- UBUNTU DEPLOY  ----------------- #
FROM ubuntu AS final
COPY --from=builder /app/server /app/server
COPY --from=builder /app/fileserver.db /app/fileserver.db
COPY --from=builder /app/idp.db /app/idp.db
COPY ./static /app/static/
COPY ./templates /app/templates/
COPY --from=builder /app/uploads /app/uploads
COPY ./secret_file /app/secret_file
RUN apt-get update -y && apt install socat sqlite3 openssl -y && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN mkdir /app/plugins
RUN mv /app/secret_file "app/uploads/REDACTED"
RUN chmod 777 /app/fileserver.db
RUN chmod 777 /app/idp.db
RUN chmod 755 /app/server

ENV DEV_FILE="REDACTED"
EXPOSE 8083

FROM pwn.red/jail
COPY --from=final / /srv
RUN echo '127.0.0.1 localhost' >> /srv/etc/hosts
COPY hook.sh /jail/hook.sh
COPY run.sh /srv/app/run
RUN chmod +x /srv/app/run /jail/hook.sh && \
    touch /srv/app/app.log && \
    chmod 777 -R /srv/app
ENV JAIL_PIDS=0 JAIL_CPU=1000 JAIL_MEM=100M JAIL_TIME=0 JAIL_ENV_DEV_FILE="REDACTED"
