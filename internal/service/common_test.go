package service

import (
	intl "github.com/s6on/mirror-docker-tags-action/internal"
	"testing"
)

func Test_platforms(t *testing.T) {
	tests := []struct {
		name string
		arg  intl.Tag
		want string
	}{
		{"", intl.Tag{Platforms: map[intl.Platform]bool{
			intl.LinuxAmd64:  true,
			intl.LinuxArmV7:  true,
			intl.LinuxMips64: true,
		}}, "linux/amd64,linux/arm/v7,linux/mips64"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := platforms(tt.arg); got != tt.want {
				t.Errorf("platforms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pullTag(t *testing.T) {
	tests := []struct {
		name string
		arg  intl.Tag
		want string
	}{
		{"", intl.Tag{Names: map[string]bool{"latest": true, "20.10": true, "groovy": true}}, "20.10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pullTag(tt.arg); got != tt.want {
				t.Errorf("pullTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toRepoName(t *testing.T) {
	type args struct {
		to   string
		repo string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Ubuntu ", args{"s6on", "ubuntu"}, "s6on/ubuntu"},
		{"Custom Repo ", args{"s6on", "user/repo"}, "s6on/user-repo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toRepoName(tt.args.to, tt.args.repo); got != tt.want {
				t.Errorf("toRepoName() = %v, want %v", got, tt.want)
			}
		})
	}
}
