# borgmatic-exporter
A metrics exporter for Borgmatic written in Go that can be scraped by Prometheus.

The exporter will output metrics on a http endpoint per repository for:

- Last backup timestamp
- Original size
- Compressed size
- Deduplicated size
- Unique chunks
- Total chunks
- Number of archives

It will also output metrics for:
- Borgmatic version
- Borg version
- Number of repositories

Example output:

```
# HELP borg_info Installed version of Borg
# TYPE borg_info gauge
borg_info{version="borg 1.1.16"} 1
# HELP borgmatic_archives_total Number of archives
# TYPE borgmatic_archives_total gauge
borgmatic_archives_total{repository="ssh://borg@server01/backup/my-data"} 13
# HELP borgmatic_chunks_total Total chunks in backup data
# TYPE borgmatic_chunks_total gauge
borgmatic_chunks_total{repository="ssh://borg@server01/backup/my-data"} 2.13976e+07
# HELP borgmatic_compressed_size Compressed size of backup data
# TYPE borgmatic_compressed_size gauge
borgmatic_compressed_size{repository="ssh://borg@server01/backup/my-data"} 5.3958861048667e+13
# HELP borgmatic_deduplicated_size Deduplicated size of backup data
# TYPE borgmatic_deduplicated_size gauge
borgmatic_deduplicated_size{repository="ssh://borg@server01/backup/my-data"} 4.251417149552e+12
# HELP borgmatic_info Installed version of Borgmatic
# TYPE borgmatic_info gauge
borgmatic_info{version="1.5.13.dev0"} 1
# HELP borgmatic_last_backup_timestamp Timestamp of latest backup
# TYPE borgmatic_last_backup_timestamp gauge
borgmatic_last_backup_timestamp{repository="ssh://borg@server01/backup/my-data"} 1.678500043e+09
# HELP borgmatic_original_size Original size of backup data
# TYPE borgmatic_original_size gauge
borgmatic_original_size{repository="ssh://borg@server01/backup/my-data"} 5.4545490460833e+13
# HELP borgmatic_repos_total Number of repositories
# TYPE borgmatic_repos_total gauge
borgmatic_repos_total 1
# HELP borgmatic_unique_chunks_total Total unique chunks in backup data
# TYPE borgmatic_unique_chunks_total gauge
borgmatic_unique_chunks_total{repository="ssh://borg@server01/backup/my-data"} 1.679772e+06
```

## Running locally
1. Download the appropriate release for your platform `curl -LO https://github.com/fishgrimsby/borgmatic-exporter/releases/download/v0.0.4/borgmatic-exporter_0.0.4_linux_amd64.tar.gz`
3. Extract the binary `tar xzvf borgmatic-exporter_0.0.4_linux_amd64.tar.gz`
4. Allow execution `chmod +x borgmatic-exporter`
5. Run the executable `sudo ./borgmatic-exporter`
6. Browse to http://localhost:8090/metrics

## Installation (Linux - systemd)
1. Download the release, extract the binary and place it at `/opt/borgmatic-exporter/borgmatic-exporter`

2. Create a file at `/etc/systemd/system/borgmatic-exporter.service` with the following contents:
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

3. Reload systemd configuration with `systemctl daemon-reload`
4. Start the service `systemctl start borgmatic-exporter.service`
5. Enable service start on boot `systemctl enable borgmatic-exporter.service`

## Docker
A multi-architecture docker image with usage instructions is available at https://hub.docker.com/r/fishgrimsby/borgmatic-exporter

## Build
1. `git clone https://github.com/fishgrimsby/borgmatic-exporter.git`
2. `cd borgmatic-exporter`
3. `mkdir bin`
4. `cd bin`
5. `go build ../cmd/borgmatic-exporter`

## Configuration options
The following configuration options are available and must be set via Environment variable:
- BORGMATIC_EXPORTER_PORT - Set the http listen port. Default to `8090`
- BORGMATIC_EXPORTER_ENDPOINT - Set the metrics endpoint. Defaults to `metrics`
- BORGMATIC_EXPORTER_CONFIG - Overrides the default config paths using `borgmatic -c`. Uses Borgmatic defaults if not set. Multiple config paths must be separated with spaces, i.e `/path/to/config1.yaml /path/to/config2.yaml`
- BORGMATIC_EXPORTER_DEBUG - Enable debug messages to stdout
- BORGMATIC_EXPORTER_LOGFORMAT - Set the format of logs. Valid options are `keyvalue` (default) and `json`