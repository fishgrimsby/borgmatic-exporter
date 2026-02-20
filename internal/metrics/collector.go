package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/fishgrimsby/borgmatic-exporter/internal/borg"
	"github.com/fishgrimsby/borgmatic-exporter/internal/borgmatic"
	"github.com/fishgrimsby/borgmatic-exporter/internal/logs"
	"github.com/prometheus/client_golang/prometheus"
)

type cachedData struct {
	borgVersion      string
	borgmaticVersion string
	listResults      []borgmatic.ListResult
	infoResults      []borgmatic.InfoResult
}

type Collector struct {
	config                     string
	timeout                    time.Duration
	interval                   time.Duration
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

	mu    sync.RWMutex
	cache *cachedData
}

func New(config string, timeout time.Duration, interval time.Duration) *Collector {

	return &Collector{
		config:                     config,
		timeout:                    timeout,
		interval:                   interval,
		borgmaticTotalUniqueChunks: prometheus.NewDesc("borgmatic_unique_chunks_total", "Total number of unique chunks in backup data", []string{"repository"}, nil),
		borgmaticTotalChunks:       prometheus.NewDesc("borgmatic_chunks_total", "Total number of chunks in backup data", []string{"repository"}, nil),
		borgmaticDeduplicatedSize:  prometheus.NewDesc("borgmatic_deduplicated_size", "Deduplicated size in bytes of backup data", []string{"repository"}, nil),
		borgmaticCompressedSize:    prometheus.NewDesc("borgmatic_compressed_size", "Compressed size in bytes of backup data", []string{"repository"}, nil),
		borgmaticOriginalSize:      prometheus.NewDesc("borgmatic_original_size", "Original size in bytes of backup data", []string{"repository"}, nil),
		borgmaticLastBackupTime:    prometheus.NewDesc("borgmatic_last_backup_timestamp", "Timestamp of latest backup", []string{"repository"}, nil),
		borgmaticRepos:             prometheus.NewDesc("borgmatic_repos_total", "Total number of repositories", nil, nil),
		borgmaticVersion:           prometheus.NewDesc("borgmatic_info", "Installed version of Borgmatic", []string{"version"}, nil),
		borgVersion:                prometheus.NewDesc("borg_info", "Installed version of Borg", []string{"version"}, nil),
		borgmaticArchives:          prometheus.NewDesc("borgmatic_archives_total", "Total number of archives", []string{"repository"}, nil),
	}
}

func (c *Collector) runCollection(ctx context.Context) {
	collectionCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	logs.Logger.Debug("Start collecting metrics")

	borgResult, err := borg.New(collectionCtx)
	if err != nil {
		logs.Logger.Error(err.Error())
		return
	}

	borgmaticResult, err := borgmatic.New(collectionCtx, c.config)
	if err != nil {
		logs.Logger.Error(err.Error())
		return
	}

	data := cachedData{
		borgVersion:      borgResult.Version,
		borgmaticVersion: borgmaticResult.Version,
		listResults:      borgmaticResult.ListResult,
		infoResults:      borgmaticResult.InfoResult,
	}

	c.mu.Lock()
	c.cache = &data
	c.mu.Unlock()

	logs.Logger.Debug("End collecting metrics")
}

func (c *Collector) Start(ctx context.Context) {
	c.runCollection(ctx)
	go func() {
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.runCollection(ctx)
			}
		}
	}()
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

func sendMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labelValues ...string) {
	m, err := prometheus.NewConstMetric(desc, valueType, value, labelValues...)
	if err != nil {
		logs.Logger.Error("failed to create metric", "error", err.Error())
		return
	}
	ch <- m
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mu.RLock()
	cache := c.cache
	c.mu.RUnlock()

	if cache == nil {
		return
	}

	sendMetric(ch, c.borgVersion, prometheus.GaugeValue, 1, cache.borgVersion)
	sendMetric(ch, c.borgmaticVersion, prometheus.GaugeValue, 1, cache.borgmaticVersion)
	sendMetric(ch, c.borgmaticRepos, prometheus.GaugeValue, float64(len(cache.listResults)))

	for _, result := range cache.listResults {
		sendMetric(ch, c.borgmaticArchives, prometheus.GaugeValue, float64(len(result.Archives)), result.Repository.Location)
		sendMetric(ch, c.borgmaticLastBackupTime, prometheus.GaugeValue, float64(borgmatic.LastBackupTime(&result)), result.Repository.Location)
	}

	for _, info := range cache.infoResults {
		sendMetric(ch, c.borgmaticDeduplicatedSize, prometheus.GaugeValue, float64(info.Cache.Stats.UniqueCsize), info.Repository.Location)
		sendMetric(ch, c.borgmaticCompressedSize, prometheus.GaugeValue, float64(info.Cache.Stats.TotalCsize), info.Repository.Location)
		sendMetric(ch, c.borgmaticOriginalSize, prometheus.GaugeValue, float64(info.Cache.Stats.TotalSize), info.Repository.Location)
		sendMetric(ch, c.borgmaticTotalChunks, prometheus.GaugeValue, float64(info.Cache.Stats.TotalChunks), info.Repository.Location)
		sendMetric(ch, c.borgmaticTotalUniqueChunks, prometheus.GaugeValue, float64(info.Cache.Stats.TotalUniqueChunks), info.Repository.Location)
	}
}
