package borg

type borg struct {
	Version string
}

func New() (*borg, error) {
	b := borg{}

	ver, err := getVersion()

	if err != nil {
		b.Version = "0"
	}

	b.Version = ver

	return &b, nil
}
