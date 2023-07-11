package internal

import (
	"time"
)

type Platform uint8

const (
	LinuxAmd64 Platform = iota
	LinuxArm64
	LinuxRiscv64
	LinuxPpc64le
	LinuxS390x
	Linux386
	LinuxMips64le
	LinuxMips64
	LinuxArmV7
	LinuxArmV6
)

func (p Platform) String() string {
	return [...]string{
		"linux/amd64",
		"linux/arm64",
		"linux/riscv64",
		"linux/ppc64le",
		"linux/s390x",
		"linux/386",
		"linux/mips64le",
		"linux/mips64",
		"linux/arm/v7",
		"linux/arm/v6",
	}[p]
}

type Tag struct {
	Names        map[string]bool
	Platforms    map[Platform]bool
	LatestUpdate time.Time
}

func NewTag() Tag {
	return Tag{Names: make(map[string]bool), Platforms: make(map[Platform]bool)}
}

type Matrix struct {
	Include []Params `json:"include"`
}

type Params struct { // TODO rename
	BaseImg   string `json:"base_img"`
	Tags      string `json:"tags"`
	Platforms string `json:"platforms"`
}
