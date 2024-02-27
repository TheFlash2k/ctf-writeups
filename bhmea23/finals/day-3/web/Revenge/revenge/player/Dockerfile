FROM ubuntu:20.04@sha256:33a5cc25d22c45900796a1aca487ad7a7cb09f09ea00b779e3b2026b4fc2faba 

ENV DEBIAN_FRONTEND=noninteractive
RUN apt update

# Setup postgresql
RUN apt install wget ca-certificates -y
RUN apt install software-properties-common -y

RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'
RUN apt update

RUN apt install postgresql postgresql-contrib -y


# Setup flask app
RUN apt install python3 -y
RUN apt install python3-pip curl unzip iptables -y


COPY share/requirement.txt /app/requirement.txt
RUN pip install -r /app/requirement.txt
EXPOSE 8080

# Setup bot
RUN pip install selenium

RUN apt install curl -y && curl -sS -o - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add
RUN bash -c "echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' >> /etc/apt/sources.list.d/google-chrome.list"
RUN apt -y update && apt -y install google-chrome-stable

WORKDIR /tmp
RUN wget https://edgedl.me.gvt1.com/edgedl/chrome/chrome-for-testing/119.0.6045.159/linux64/chromedriver-linux64.zip
RUN unzip chromedriver-linux64.zip
RUN mv chromedriver-linux64/chromedriver /usr/bin/chromedriver && chmod 555 /usr/bin/chromedriver
WORKDIR /

# Setup users
RUN useradd -m flaskuser

# Setup nginx
COPY nginx_1.16.1.deb /nginx.deb
RUN dpkg -i /nginx.deb
RUN rm /nginx.deb
#RUN apt install -y nginx
COPY ./nginx/default.conf /etc/nginx/conf.d/nginx.conf
COPY ./nginx/www/ /usr/share/nginx/html/
RUN chmod 555 /usr/share/nginx/html/ -R
RUN rm -f /usr/share/nginx/html/50x.html

# Python webapp
COPY pyserve /pyserve
RUN chmod 777 /pyserve -R

#init db
COPY init.sql /init.sql
#cpy project
COPY share /share
# Setup entrypoint
COPY entrypoint.sh /
RUN chmod 555 /entrypoint.sh

# Start the server
CMD /entrypoint.sh

