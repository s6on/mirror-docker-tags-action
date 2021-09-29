package main

import (
	"reflect"
	"testing"
)

func Test_getRepos(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want map[string][]string
	}{
		{"", "ubuntu[latest,20.10,18.04],alpine[latest,3]", map[string][]string{
			"ubuntu": {"latest", "20.10", "18.04"},
			"alpine": {"latest", "3"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRepos(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRepos() = %v, want %v", got, tt.want)
			}
		})
	}
}
