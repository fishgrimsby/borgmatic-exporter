# borgmatic‑exporter

A Prometheus exporter written in Go that exposes metrics collected from a Borgmatic installation.

---

## Metrics
The exporter provides the following metrics per repository:

| Metric | Description |
|--------|-------------|
| `borgmatic_last_backup_timestamp` | Unix timestamp of the latest backup |
| `borgmatic_original_size` | Total size of the original data |
| `borgmatic_compressed_size` | Size after compression |
| `borgmatic_deduplicated_size` | Deduplicated size |
| `borgmatic_unique_chunks_total` | Number of unique chunks |
| `borgmatic_chunks_total` | Total number of chunks |
| `borgmatic_archives_total` | Number of archives |

In addition, the following global metrics are exported:

| Metric | Description |
|--------|-------------|
| `borgmatic_info` | Installed Borgmatic version |
| `borg_info` | Installed Borg version |
| `borgmatic_repos_total` | Number of repositories configured in Borgmatic |

Example output (Prometheus text format):
```
# HELP borg_info Installed version of Borg
# TYPE borg_info gauge
borg_info{version="borg 1.1.16"} 1
# HELP borgmatic_archives_total Number of archives
# TYPE borgmatic_archives_total gauge
borgmatic_archives_total{repository="ssh://borg@server01/backup/my-data"} 13
# HELP borgmatic_last_backup_timestamp Timestamp of latest backup
# TYPE borgmatic_last_backup_timestamp gauge
borgmatic_last_backup_timestamp{repository="ssh://borg@server01/backup/my-data"} 1.678500043e+09
... (other metrics omitted) ...
```

---

## Prerequisites
* Go 1.23 or later to build from source.
* `borg` and `borgmatic` binaries must be in the `PATH` as the exporter invokes them directly.
* Environment variables (see **Configuration** section) for runtime settings.

---

## Building from Source
```bash
# Clone the repository
git clone https://github.com/fishgrimsby/borgmatic-exporter.git
cd borgmatic-exporter

# Build the exporter binary
go build ./cmd/borgmatic-exporter

# Run the binary (see Running section for options)
./borgmatic-exporter
```

---

## Downloading Pre-Built Binaries
Binaries are availble for multiple platforms at https://github.com/fishgrimsby/borgmatic-exporter/releases

---

## Running the Exporter
The exporter listens on `<HOST>:<PORT>/<ENDPOINT>`.

### Default (no env vars)
```bash
# Default settings
# HOST=0.0.0.0, PORT=8090, ENDPOINT=metrics
./borgmatic-exporter
# Browse to http://localhost:8090/metrics
```

### With custom configuration
Set environment variables before launching:
```bash
export BORGMATIC_EXPORTER_HOST=0.0.0.0
export BORGMATIC_EXPORTER_PORT=9090
export BORGMATIC_EXPORTER_ENDPOINT=metrics
export BORGMATIC_EXPORTER_CONFIG="/etc/borgmatic/backup1.yaml /etc/borgmatic/backup2.yaml"
export BORGMATIC_EXPORTER_DEBUG=true
export BORGMATIC_EXPORTER_LOGFORMAT=json
export BORGMATIC_EXPORTER_TIMEOUT=30s

./borgmatic-exporter
```

---

## Systemd Service (Linux)
Create a service unit at `/etc/systemd/system/borgmatic-exporter.service`:
```
[Unit]
Description=Borgmatic Exporter Daemon
After=syslog.target network.target

[Service]
User=root
Group=root
Type=simple
ExecStart=/opt/borgmatic-exporter/borgmatic-exporter
TimeoutStopSec=20
KillMode=process
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
Then enable and start:
```bash
systemctl daemon-reload
systemctl enable --now borgmatic-exporter.service
```

---

## Docker
A multi-architecture docker image with usage instructions is available at https://hub.docker.com/r/fishgrimsby/borgmatic-exporter

---

## Configuration Options
The exporter reads the following environment variables (prefixed with `BORGMATIC_EXPORTER_`). If unset, defaults apply.

| Variable | Default | Description |
|----------|---------|-------------|
| `HOST` | `0.0.0.0` | HTTP listen address |
| `PORT` | `8090` | HTTP listen port |
| `ENDPOINT` | `metrics` | Path served by Prometheus |
| `CONFIG` | `""` | Space‑separated list of Borgmatic config files (`borgmatic -c <config>`). If omitted, Borgmatic defaults are used |
| `DEBUG` | `false` | Enable debug‑level logging |
| `LOGFORMAT` | `keyvalue` | Log format (`keyvalue` or `json`) |
| `TIMEOUT` | `120s` | Timeout for querying metric values |

---

## Testing
Run all unit tests with:
```bash
go test ./...
```
Run a specific test:
```bash
go test -run TestName ./internal/borgmatic
```

---

## Contribution Guidelines
Pull requests are welcome. Please follow the existing code style, run `go fmt`, and ensure all tests pass.

---

## License
MIT
