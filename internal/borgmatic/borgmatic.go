package borgmatic

import (
	"context"
	"sort"
	"time"
)

type borgmatic struct {
	Version     string
	ListResult  []ListResult
	InfoResult []InfoResult
}

func New(ctx context.Context, config string) (*borgmatic, error) {
	b := borgmatic{}

	// Get Version
	ver, err := getVersion(ctx)

	if err != nil {
		b.Version = "0"
	} else {
		b.Version = ver
	}

	// Get List Results
	res, err := getArchives(ctx, config)

	if err != nil {
		res = []ListResult{}
	}

	b.ListResult = res

	// Get Info Results
	info, err := getInfo(ctx, config)

	if err != nil {
		info = []InfoResult{}
	}

	b.InfoResult = info

	return &b, nil
}

func (b *borgmatic) LastBackupTime(result *ListResult) int64 {
	var lastBackupTime int64 = 0

	var times []string
	for _, result := range result.Archives {
		times = append(times, result.Time)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(times)))

	if len(times) > 0 {
		convert, err := time.Parse("2006-01-02T15:04:05.999999999", times[0])

		if err != nil {
			lastBackupTime = 0
		} else {
			lastBackupTime = convert.Unix()
		}

	}

	return int64(lastBackupTime)
}
