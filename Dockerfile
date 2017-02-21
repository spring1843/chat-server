FROM ubuntu:wily
#Get Ubuntu

ADD https://github.com/spring1843/chat-server/releases/download/v0.0.1/latest_package.zip /tmp/latest_package.zip

RUN apt-get update -y && apt-get install --no-install-recommends -y -q zip unzip

RUN unzip /tmp/latest_package.zip -d /tmp

EXPOSE 80
ENTRYPOINT ["/tmp/chat-server.linux.amd64", "-config", "/tmp/config.json"]
