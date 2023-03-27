package borgmatic

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

var infoResult InfoResult = InfoResult{
	Cache:       infoCache,
	Encryption:  infoEncryption,
	Repository:  infoRepository,
	SecurityDir: "/Users/user1/.config/borg/security/c517781f67da49cf8396a29227e4ad68591cde71baa0e04acc18893608f869de",
}
var infoCache Cache = Cache{
	Path:  "/Users/user1/.cache/borg/c517781f67da49cf8396a29227e4ad68591cde71baa0e04acc18893608f869de",
	Stats: infoStats,
}
var infoRepository Repository = Repository{
	Id:           "c517781f67da49cf8396a29227e4ad68591cde71baa0e04acc18893608f869de",
	LastModified: "2023-03-15T20:57:35.000000",
	Location:     "/Users/user1/borgmatic/backup1",
}

var infoStats Stats = Stats{
	TotalChunks:       45437,
	TotalCsize:        4619223355,
	TotalSize:         8385768055,
	TotalUniqueChunks: 9697,
	UniqueCsize:       2007533768,
	UniqueSize:        3279757616,
}

var infoEncryption Encryption = Encryption{
	Mode: "repokey",
}

func TestGetInfoSingleRepository(t *testing.T) {
	execCommand = fakeExecCommandGetInfoSingleRepository
	defer func() { execCommand = exec.Command }()

	want := []InfoResult{infoResult}
	got, err := getInfo("")

	if err != nil {
		t.Fatalf("Expected nil error, got %#v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf(`getInfo(config) = %v, want match for %#v, nil`, got, want)
	}
}

func fakeExecCommandGetInfoSingleRepository(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestGetInfoSingleRepositoryHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

const borgmaticGetInfoResult = `[{"cache": {"path": "/Users/user1/.cache/borg/c517781f67da49cf8396a29227e4ad68591cde71baa0e04acc18893608f869de", "stats": {"total_chunks": 45437, "total_csize": 4619223355, "total_size": 8385768055, "total_unique_chunks": 9697, "unique_csize": 2007533768, "unique_size": 3279757616}}, "encryption": {"mode": "repokey"}, "repository": {"id": "c517781f67da49cf8396a29227e4ad68591cde71baa0e04acc18893608f869de", "last_modified": "2023-03-15T20:57:35.000000", "location": "/Users/user1/borgmatic/backup1"}, "security_dir": "/Users/user1/.config/borg/security/c517781f67da49cf8396a29227e4ad68591cde71baa0e04acc18893608f869de"}]`

func TestGetInfoSingleRepositoryHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprintf(os.Stdout, borgmaticGetInfoResult)
	os.Exit(0)
}
