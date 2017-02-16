FROM ubuntu:wily
#Get Ubuntu

# Add linux binary
ADD ./chat-server.linux.amd64 /tmp/chat-server.linux.amd64

# Add config.json
ADD ./config.json /tmp/config.json

# add static web files
ADD ./static-web /tmp/static-web

# Add certificats
ADD ./tls.crt /tmp/tls.crt
ADD ./tls.key /tmp/tls.key

EXPOSE 80
ENTRYPOINT ["/tmp/chat-server.linux.amd64", "-config", "/tmp/config.json"]
