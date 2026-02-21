package borgmatic

import (
	"context"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	config := ""
	want := &borgmatic{}
	got, err := New(context.Background(), config)

	if err != nil {
		t.Fatalf("Expected nil error, got %#v", err)
	}

	if reflect.TypeOf(want) != reflect.TypeOf(got) {
		t.Fatalf(`New() = %v, want match for %#q, nil`, reflect.TypeOf(got), reflect.TypeOf(want))
	}
}

func TestLastBackupTime(t *testing.T) {
	var want int64 = 1678913859
	got := LastBackupTime(&listResult)

	if want != got {
		t.Fatalf(`LastBackupTime(listResult) = %v, want match for %#v`, got, want)
	}
}

func TestLastBackupTimeZeroArchives(t *testing.T) {
	listResultZeroArchives := listResult
	listResultZeroArchives.Archives = []Archive{}
	var want int64 = 0
	got := LastBackupTime(&listResultZeroArchives)

	if want != got {
		t.Fatalf(`LastBackupTime(listResultZeroArchives) = %v, want match for %#v`, got, want)
	}
}
