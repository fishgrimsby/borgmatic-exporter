package borg

import "context"

type borg struct {
	Version string
}

func New(ctx context.Context) (*borg, error) {
	b := borg{}

	ver, err := getVersion(ctx)

	if err != nil {
		b.Version = "0"
	} else {
		b.Version = ver
	}

	return &b, nil
}
