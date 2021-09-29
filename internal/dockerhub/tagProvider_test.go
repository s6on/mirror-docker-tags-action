package dockerhub

import (
	intl "github.com/s6on/mirror-docker-tags-action/internal"
	dckr "github.com/s6on/mirror-docker-tags-action/lib/dockerhub"
	"testing"
)

func TestGet(t *testing.T) {
	type args struct {
		repo string
		tags []string
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{"Ubuntu 2", args{repo: "ubuntu", tags: []string{"20.10", "foo", "20.04"}}, 2},
		{"Alpine 1", args{repo: "alpine", tags: []string{"latest"}}, 1},
		{"Debian 0", args{repo: "debian", tags: []string{"foo"}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.repo, tt.args.tags); len(got) != int(tt.want) {
				t.Errorf("Get() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_platform(t *testing.T) {
	tests := []struct {
		name  string
		arg   dckr.Image
		want  intl.Platform
		want1 bool
	}{
		{"windows", dckr.Image{Os: "windows"}, 0, false},
		{"amd64", dckr.Image{Os: "linux", Architecture: "amd64"}, intl.LinuxAmd64, true},
		{"arm64", dckr.Image{Os: "linux", Architecture: "arm", Variant: "v8"}, intl.LinuxArm64, true},
		{"armhf", dckr.Image{Os: "linux", Architecture: "arm", Variant: "v6"}, intl.LinuxArmV6, true},
		{"arm_unknown", dckr.Image{Os: "linux", Architecture: "arm", Variant: "v2"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := platform(tt.arg)
			if got != tt.want {
				t.Errorf("platform() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("platform() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
