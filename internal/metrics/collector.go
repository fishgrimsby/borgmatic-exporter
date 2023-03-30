package metrics

import (
	"github.com/fishgrimsby/borgmatic-exporter/internal/borg"
	"github.com/fishgrimsby/borgmatic-exporter/internal/borgmatic"
	"github.com/fishgrimsby/borgmatic-exporter/internal/logs"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	config                     string
	borgmaticTotalUniqueChunks *prometheus.Desc
	borgmaticTotalChunks       *prometheus.Desc
	borgmaticDeduplicatedSize  *prometheus.Desc
	borgmaticCompressedSize    *prometheus.Desc
	borgmaticOriginalSize      *prometheus.Desc
	borgmaticLastBackupTime    *prometheus.Desc
	borgmaticRepos             *prometheus.Desc
	borgmaticArchives          *prometheus.Desc
	borgmaticVersion           *prometheus.Desc
	borgVersion                *prometheus.Desc
}

func New(config string) *Collector {

	return &Collector{
		config:                     config,
		borgmaticTotalUniqueChunks: prometheus.NewDesc("borgmatic_unique_chunks_total", "Total number of unique chunks in backup data", []string{"repository"}, nil),
		borgmaticTotalChunks:       prometheus.NewDesc("borgmatic_chunks_total", "Total number of chunks in backup data", []string{"repository"}, nil),
		borgmaticDeduplicatedSize:  prometheus.NewDesc("borgmatic_deduplicated_size", "Deduplicated size in bytes of backup data", []string{"repository"}, nil),
		borgmaticCompressedSize:    prometheus.NewDesc("borgmatic_compressed_size", "Compressed size in bytes of backup data", []string{"repository"}, nil),
		borgmaticOriginalSize:      prometheus.NewDesc("borgmatic_original_size", "Original size in bytes of backup data", []string{"repository"}, nil),
		borgmaticLastBackupTime:    prometheus.NewDesc("borgmatic_last_backup_timestamp", "Timestamp of latest backup", []string{"repository"}, nil),
		borgmaticRepos:             prometheus.NewDesc("borgmatic_repos_total", "Total number of repositories", nil, nil),
		borgmaticVersion:           prometheus.NewDesc("borgmatic_info", "Installed version of Borgmatic", []string{"version"}, nil),
		borgVersion:                prometheus.NewDesc("borg_info", "Installed version of Borg", []string{"version"}, nil),
		borgmaticArchives:          prometheus.NewDesc("borgmatic_archives_total", "Total umber of archives", []string{"repository"}, nil),
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.borgmaticTotalUniqueChunks
	ch <- collector.borgmaticTotalChunks
	ch <- collector.borgmaticDeduplicatedSize
	ch <- collector.borgmaticCompressedSize
	ch <- collector.borgmaticOriginalSize
	ch <- collector.borgmaticLastBackupTime
	ch <- collector.borgmaticRepos
	ch <- collector.borgmaticArchives
	ch <- collector.borgmaticVersion
	ch <- collector.borgVersion
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	logs.Logger.Debug("Start collecting metrics")
	borg, err := borg.New()
	if err != nil {
		logs.Logger.Error(err.Error())
	}

	ch <- prometheus.MustNewConstMetric(collector.borgVersion, prometheus.GaugeValue, 1, borg.Version)

	borgmatic, err := borgmatic.New(collector.config)
	if err != nil {
		logs.Logger.Error(err.Error())
	}
	ch <- prometheus.MustNewConstMetric(collector.borgmaticVersion, prometheus.GaugeValue, 1, borgmatic.Version)
	ch <- prometheus.MustNewConstMetric(collector.borgmaticRepos, prometheus.GaugeValue, float64(len(borgmatic.ListResult)))

	for _, result := range borgmatic.ListResult {
		ch <- prometheus.MustNewConstMetric(collector.borgmaticArchives, prometheus.GaugeValue, float64(len(result.Archives)), result.Repository.Location)
		ch <- prometheus.MustNewConstMetric(collector.borgmaticLastBackupTime, prometheus.GaugeValue, float64(borgmatic.LastBackupTime(&result)), result.Repository.Location)
	}

	for _, info := range borgmatic.InfoResullt {
		ch <- prometheus.MustNewConstMetric(collector.borgmaticDeduplicatedSize, prometheus.GaugeValue, float64(info.Cache.Stats.UniqueCsize), info.Repository.Location)
		ch <- prometheus.MustNewConstMetric(collector.borgmaticCompressedSize, prometheus.GaugeValue, float64(info.Cache.Stats.TotalCsize), info.Repository.Location)
		ch <- prometheus.MustNewConstMetric(collector.borgmaticOriginalSize, prometheus.GaugeValue, float64(info.Cache.Stats.TotalSize), info.Repository.Location)
		ch <- prometheus.MustNewConstMetric(collector.borgmaticTotalChunks, prometheus.GaugeValue, float64(info.Cache.Stats.TotalChunks), info.Repository.Location)
		ch <- prometheus.MustNewConstMetric(collector.borgmaticTotalUniqueChunks, prometheus.GaugeValue, float64(info.Cache.Stats.TotalUniqueChunks), info.Repository.Location)
	}

	logs.Logger.Debug("End collecting metrics")
}
