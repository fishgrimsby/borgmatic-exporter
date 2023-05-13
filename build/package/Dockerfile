FROM golang:1.20 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build ./cmd/borgmatic-exporter

FROM build-stage AS test-stage
RUN go test -v ./...

FROM debian:bullseye-slim AS package-stage
WORKDIR /
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
    && apt-get -y upgrade \
    && apt-get -y install --no-install-recommends \
    borgbackup borgmatic openssh-client \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

FROM package-stage AS final-stage
WORKDIR /
COPY --from=build-stage /app/borgmatic-exporter /bin/borgmatic-exporter
EXPOSE 8090
ENTRYPOINT ["/bin/borgmatic-exporter"]
