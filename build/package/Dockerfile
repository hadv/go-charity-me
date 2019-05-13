FROM golang:1.12-alpine AS build
RUN apk --no-cache add git build-base curl

WORKDIR /go/bin
ENV MIGRATE_VERSION="v4.2.5"
RUN  curl -L https://github.com/golang-migrate/migrate/releases/download/$MIGRATE_VERSION/migrate.linux-amd64.tar.gz | tar xvz
ENV WAITFORIT_VERSION="v2.4.1"
RUN curl -o waitforit -sSL https://github.com/maxcnunes/waitforit/releases/download/$WAITFORIT_VERSION/waitforit-linux_amd64 && \
    chmod +x waitforit

ADD . /src
WORKDIR /src
RUN go install ./cmd/charity

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
COPY --from=build /src/db /db
COPY --from=build /go/bin/* /bin/
CMD ["/bin/charity"]
