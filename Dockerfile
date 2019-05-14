FROM golang:1.12

RUN mkdir -p $GOPATH/src/github.com/psprings/code-concierge

WORKDIR $GOPATH/src/github.com/psprings/code-concierge

COPY . .

RUN go get -v ./... && \
    CGO_ENABLED=0 go build -v -ldflags="-s -w" -o /code-concierge

FROM codercom/code-server

USER root

COPY docker-entrypoint.sh /home/coder/docker-entrypoint.sh

RUN echo 'coder ALL= NOPASSWD: /usr/bin/apt-get' >> /etc/sudoers && \
    chmod +x /home/coder/docker-entrypoint.sh

USER coder

COPY --from=0 /code-concierge /usr/local/bin/code-concierge

ENTRYPOINT [ "/home/coder/docker-entrypoint.sh" ]