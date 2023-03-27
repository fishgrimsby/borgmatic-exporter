package borgmatic

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

var archives []Archive = []Archive{
	{
		Archive:  "Machine-1-2023-03-11T19:45:06.169335",
		Barchive: "Machine-1-2023-03-11T19:45:06.169335 ",
		Id:       "b243952402483f584a6d7a9c50fe672f7ce2004d635f9ebb29aed78236ae7945",
		Name:     "Machine-1-2023-03-11T19:45:06.169335",
		Start:    "2023-03-11T19:45:06.000000",
		Time:     "2023-03-11T19:45:06.000000",
	},
	{
		Archive:  "Machine-1-2023-03-15T20:57:39.075581",
		Barchive: "Machine-1-2023-03-15T20:57:39.075581",
		Id:       "4b596553c24cb1e5e709cfc95d77580fc8774005ad229cc39ce937539886a749",
		Name:     "Machine-1-2023-03-15T20:57:39.075581",
		Start:    "2023-03-15T20:57:39.000000",
		Time:     "2023-03-15T20:57:39.000000",
	},
}

var encryption Encryption = Encryption{
	Mode: "repokey",
}

var repository Repository = Repository{
	Id:           "c517781f67da49cf8396a29227e4ad68591cde71baa0e04acc18893608f869de",
	LastModified: "2023-03-15T20:57:35.000000",
	Location:     "/path/to/backup/repository",
}

var listResult ListResult = ListResult{
	Archives:   archives,
	Encryption: encryption,
	Repository: repository,
}

func TestGetArchivesSingleRepository(t *testing.T) {
	execCommand = fakeExecCommandGetArchivesSingleRepository
	defer func() { execCommand = exec.Command }()

	want := []ListResult{listResult}
	got, err := getArchives("")

	if err != nil {
		t.Fatalf("Expected nil error, got %#v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf(`getArchives(config) = %v, want match for %#v, nil`, got, want)
	}

}

func fakeExecCommandGetArchivesSingleRepository(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestGetArchivesSingleRepositoryHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

const borgmaticGetArchivesSingleRepositoryResult = `[{"archives":[{"archive":"Machine-1-2023-03-11T19:45:06.169335","barchive":"Machine-1-2023-03-11T19:45:06.169335 ","id":"b243952402483f584a6d7a9c50fe672f7ce2004d635f9ebb29aed78236ae7945","name":"Machine-1-2023-03-11T19:45:06.169335","start":"2023-03-11T19:45:06.000000","time":"2023-03-11T19:45:06.000000"},{"archive":"Machine-1-2023-03-15T20:57:39.075581","barchive":"Machine-1-2023-03-15T20:57:39.075581","id":"4b596553c24cb1e5e709cfc95d77580fc8774005ad229cc39ce937539886a749","name":"Machine-1-2023-03-15T20:57:39.075581","start":"2023-03-15T20:57:39.000000","time":"2023-03-15T20:57:39.000000"}],"encryption":{"mode":"repokey"},"repository":{"id":"c517781f67da49cf8396a29227e4ad68591cde71baa0e04acc18893608f869de","last_modified":"2023-03-15T20:57:35.000000","location":"/path/to/backup/repository"}}]`

func TestGetArchivesSingleRepositoryHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprintf(os.Stdout, borgmaticGetArchivesSingleRepositoryResult)
	os.Exit(0)
}
