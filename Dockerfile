FROM golang:1.18.6-alpine3.16 AS builder

ARG USER=docker
ARG home=/code

WORKDIR $home

RUN addgroup $USER && \
    adduser -D -G $USER -h $home $USER && \
    chown -R $USER:$USER $home

USER $USER


# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

EXPOSE 8080
CMD ["app"]