# build
FROM golang:1.17.5-stretch AS build-env

WORKDIR /src

ADD . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

# run
FROM alpine:3.15

RUN apk add --no-cache \
    ca-certificates

RUN addgroup -S app \
    && adduser -S -g app app

WORKDIR /home/app

COPY --from=build-env /src/app .

RUN chown -R app:app ./

USER app
ENV USER=app

ENTRYPOINT ["./app"]
