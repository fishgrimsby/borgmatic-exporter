FROM golang:1.26.0-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/borgmatic-exporter

FROM build AS test
RUN go test -v ./...
RUN go vet -v ./...

FROM debian:trixie-slim AS final
WORKDIR /
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
 && apt-get -y upgrade \
 && apt-get -y install --no-install-recommends \
    borgbackup borgmatic openssh-client \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/borgmatic-exporter /bin/borgmatic-exporter

EXPOSE 8090
ENTRYPOINT ["/bin/borgmatic-exporter"]
