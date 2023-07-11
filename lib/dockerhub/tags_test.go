package dockerhub

import (
	"testing"
	"time"
)

func TestGetTags(t *testing.T) {
	type args struct {
		repo     string
		pageSize uint
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{"Ubuntu 25", args{repo: "ubuntu", pageSize: 25}, 25},
		{"Alpine 10", args{repo: "alpine", pageSize: 10}, 10},
		{"Debian 100", args{repo: "debian", pageSize: 100}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTags(tt.args.repo, tt.args.pageSize); len(got.Tags) != int(tt.want) {
				t.Errorf("GetTags() = %v, want %v", len(got.Tags), tt.want)
			}
		})
	}
}

func TestTag_Id(t1 *testing.T) {
	tests := []struct {
		name string
		tag  Tag
		want string
	}{
		{"Unique tags", Tag{Images: []Image{
			{Digest: "c"},
			{Digest: "a"},
			{Digest: "b"},
		}}, "a,b,c"},
		{"Repeated tags", Tag{Images: []Image{
			{Digest: "b"},
			{Digest: "b"},
			{Digest: "a"},
		}}, "a,b"},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if got := tt.tag.Id(); got != tt.want {
				t1.Errorf("Id() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActive_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want Active
	}{
		{"Active", "\"active\"", true},
		{"Null", "null", false},
		{"Empty", "", false},
		{"Inactive", "inactive", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			active := Active(false)
			if err := active.UnmarshalJSON([]byte(tt.arg)); (err != nil) || active != tt.want {
				t.Errorf("UnmarshalJSON(%v) error = %v, got = %v, want %v", tt.arg, err, active, tt.want)
			}
		})
	}
}

func TestTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want Time
	}{
		{"Null", "null", Time{}},
		{"Empty", "", Time{}},
		{"Parsing", "\"2019-06-08T17:34:35.543678\"", Time(time.Date(2019, time.June, 8, 17, 34, 35, 540000000, time.UTC))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tim := Time{}
			if err := tim.UnmarshalJSON([]byte(tt.arg)); err != nil || tim != tt.want {
				t.Errorf("UnmarshalJSON(%v) error = %v, got = %v, want %v", tt.arg, err, time.Time(tim).String(), time.Time(tt.want).String())
			}
		})
	}
}

func Test_cleanup(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"Official image without path", "ubuntu", "library/ubuntu"},
		{"With underscore", "_/ubuntu", "library/ubuntu"},
		{"Official as is", "library/alpine", "library/alpine"},
		{"As is", "s6on/ubuntu", "s6on/ubuntu"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanup(tt.arg); got != tt.want {
				t.Errorf("cleanup(%v) = %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}
