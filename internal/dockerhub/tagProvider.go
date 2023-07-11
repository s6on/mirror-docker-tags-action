package dockerhub

import (
	"log"
	"time"

	intl "github.com/s6on/mirror-docker-tags-action/internal"
	dckr "github.com/s6on/mirror-docker-tags-action/lib/dockerhub"
)

func Get(repo string, tags []string) []intl.Tag {
	rs := dckr.GetTags(repo, 100)
	tm := newTagManager()
	tm.add(rs.Tags)

	for len(rs.Next) > 0 && !tm.contains(tags) {
		rs = dckr.Get(rs.Next)
		tm.add(rs.Tags)
	}

	return tm.get(tags)
}

type tagManager struct {
	id   map[string]intl.Tag
	name map[string]intl.Tag
}

func newTagManager() tagManager {
	return tagManager{
		id:   make(map[string]intl.Tag),
		name: make(map[string]intl.Tag),
	}
}

func (tm tagManager) add(dTags []dckr.Tag) {
	for _, dt := range dTags {
		if !dt.Active {
			continue
		}
		id := dt.Id()
		tag, exist := tm.id[id]
		if !exist {
			tag = intl.NewTag()
			tm.id[id] = tag
		}
		populate(&tag, dt)
		tm.name[dt.Name] = tag
	}
}

func (tm tagManager) contains(names []string) bool {
	for _, name := range names {
		if _, exist := tm.name[name]; !exist {
			return false
		}
	}
	return true
}

func (tm tagManager) get(tn []string) []intl.Tag {
	tags := make([]intl.Tag, 0)
	nm := make(map[string]bool)
	for _, name := range tn {
		nm[name] = true
	}
	for _, n := range tn {
		if !nm[n] {
			continue
		}
		tag, exist := tm.name[n]
		if !exist {
			continue
		}
		tags = append(tags, tag)
		for name := range tag.Names {
			delete(nm, name)
		}
	}
	return tags
}

func populate(t *intl.Tag, dt dckr.Tag) {
	t.Names[dt.Name] = true
	dtt := time.Time(dt.LastUpdated)
	if dtt.After(t.LatestUpdate) {
		t.LatestUpdate = dtt
	}
	for _, img := range dt.Images {
		if !img.Active {
			continue
		}
		if p, ok := platform(img); ok {
			t.Platforms[p] = true
		}
	}
}

func platform(img dckr.Image) (intl.Platform, bool) {
	if img.Os != "linux" {
		return 0, false
	}
	switch img.Architecture {
	case "amd64":
		return intl.LinuxAmd64, true
	case "arm64":
		return intl.LinuxArm64, true
	case "riscv64":
		return intl.LinuxRiscv64, true
	case "ppc64le":
		return intl.LinuxPpc64le, true
	case "s390x":
		return intl.LinuxS390x, true
	case "386":
		return intl.Linux386, true
	case "mips64le":
		return intl.LinuxMips64le, true
	case "mips64":
		return intl.LinuxMips64, true
	case "arm":
		switch img.Variant {
		case "v7":
			return intl.LinuxArmV7, true
		case "v6":
			return intl.LinuxArmV6, true
		case "v5":
			return intl.LinuxArmV6, true
		case "v8":
			return intl.LinuxArm64, true
		default:
			log.Printf("Unknown arm variant %s", img.Variant)
		}
	default:
		log.Printf("Unknown architecture %s", img.Architecture)
	}
	return 0, false
}
