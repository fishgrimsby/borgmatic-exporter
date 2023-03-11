package borgmatic

import (
	"sort"
	"time"
)

type borgmatic struct {
	Version     string
	ListResult  []ListResult
	InfoResullt []InfoResult
}

func New(config string) (*borgmatic, error) {
	b := borgmatic{}

	// Get Version
	ver, err := getVersion()

	if err != nil {
		b.Version = "0"
	}

	b.Version = ver

	// Get List Results
	res, err := getArchives(config)

	if err != nil {
		res = []ListResult{}
	}

	b.ListResult = res

	// Get Info Results
	info, err := getInfo(config)

	if err != nil {
		info = []InfoResult{}
	}

	b.InfoResullt = info

	return &b, nil
}

func (b *borgmatic) LastBackupTime(result *ListResult) int64 {
	var lastBackupTime int64 = 0

	var times []string
	for _, result := range result.Archives {
		times = append(times, result.Time)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(times)))

	convert, err := time.Parse("2006-01-02T15:04:05.999999999", times[0])

	if err != nil {
		lastBackupTime = 0
	}

	lastBackupTime = convert.Unix()

	return int64(lastBackupTime)
}
