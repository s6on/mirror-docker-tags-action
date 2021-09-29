FROM golang:alpine

WORKDIR /go/src/mirror-docker-tags-action
COPY . .
RUN go install -v ./...

ENTRYPOINT ["action"]