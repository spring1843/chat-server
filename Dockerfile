FROM ubuntu:wily

RUN apt-get update -y && apt-get install --no-install-recommends -y -q curl build-essential ca-certificates git mercurial bzr
RUN mkdir /goroot && curl https://storage.googleapis.com/golang/go1.5.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1

RUN mkdir /gopath

ENV GOROOT /goroot
ENV GOPATH /gopath
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

RUN mkdir -p  $GOPATH/src/github.com/spring1843/chat-server
WORKDIR  $GOPATH/src/github.com/spring1843/chat-server
ADD . $GOPATH/src/github.com/spring1843/chat-server

EXPOSE 4000 4001 4004
CMD ["/gopath/src/github.com/spring1843/chat-server/src/srceb_run.sh"]