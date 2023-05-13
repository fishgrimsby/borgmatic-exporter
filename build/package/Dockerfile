FROM --platform=$BUILDPLATFORM golang:1.20 AS build-stage
ARG BUILDPLATFORM
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build ./cmd/borgmatic-exporter
RUN ls -al

FROM build-stage AS test-stage
RUN go test -v ./...

FROM --platform=$BUILDPLATFORM debian:bullseye-slim AS package-stage
ARG BUILDPLATFORM
WORKDIR /
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
    && apt-get -y upgrade \
    && apt-get -y install --no-install-recommends \
    borgbackup borgmatic \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

FROM --platform=$BUILDPLATFORM package-stage AS final-stage
ARG BUILDPLATFORM
WORKDIR /
COPY --from=build-stage /app/borgmatic-exporter /bin/borgmatic-exporter
EXPOSE 8090
ENTRYPOINT ["/bin/borgmatic-exporter"]
