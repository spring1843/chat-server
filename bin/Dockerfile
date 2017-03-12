FROM ubuntu:wily
#Get Ubuntu

ADD . /tmp/chat-server/

EXPOSE 80
ENTRYPOINT ["/tmp/chat-server/chat-server.linux.amd64", "-config", "/tmp/chat-server/config.json"]
