# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git mercurial gcc
ADD . /src
RUN cd /src && go mod vendor && go build -o goapp

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/goapp /app/
EXPOSE 8080
ENTRYPOINT ./goapp