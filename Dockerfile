# build stage
FROM golang:1.9.1-alpine AS build-env

RUN apk update && apk upgrade && apk add --no-cache bash git

ENV GOPATH /go
ENV SOURCES /go/src/go-microservice/

RUN go get github.com/gin-gonic/gin; \
 go get github.com/jmoiron/sqlx; \
 go get github.com/lib/pq; \
 go get golang.org/x/crypto/bcrypt

COPY gin-framework/ ${SOURCES}/gin-framework/

RUN cd ${SOURCES}gin-framework && CGO_ENABLED=0 go build -o goapp


# final stage
FROM alpine
WORKDIR /app
ENV SOURCES /go/src/go-microservice/
COPY --from=build-env ${SOURCES}/gin-framework/ /app/gin-framework/

ENTRYPOINT cd gin-framework && ./goapp