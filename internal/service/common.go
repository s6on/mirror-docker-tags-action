package service

import (
	"sort"
	"strings"

	intl "github.com/s6on/mirror-docker-tags-action/internal"
	dckr "github.com/s6on/mirror-docker-tags-action/internal/dockerhub"
)

type MatrixBuilder struct {
	From             map[string][]string
	To               string
	ExtraRegistry    string
	UpdateAll        bool
	AllowedPlatforms map[string]bool
}

func (m *MatrixBuilder) Get() intl.Matrix {
	parms := make([]intl.Params, 0)
	for repo, tags := range m.From {
		ot := dckr.Get(repo, tags)
		trn := m.toRepoName(repo)
		var nt []intl.Tag
		if m.UpdateAll {
			nt = ot
		} else {
			tt := dckr.Get(trn, tags)
			nt = getTagsToUpdate(ot, tt)
		}
		for _, tag := range nt {
			parms = append(parms, m.param(tag, repo, trn))
		}
	}
	return intl.Matrix{Include: parms}
}

func (m *MatrixBuilder) toRepoName(repo string) string {
	return m.To + "/" + strings.ReplaceAll(repo, "/", "-")
}

func (m *MatrixBuilder) param(tag intl.Tag, fromRepo string, toRepo string) intl.Params {
	return intl.Params{
		BaseImg:   fromRepo + ":" + pullTag(tag),
		Tags:      m.tags(tag, toRepo),
		Platforms: m.platforms(tag),
	}
}

func (m *MatrixBuilder) tags(tag intl.Tag, toRepo string) string {
	t := make([]string, 0)
	for name := range tag.Names {
		t = append(t, toRepo+":"+name)
		if len(m.ExtraRegistry) > 0 {
			t = append(t, m.ExtraRegistry+"/"+toRepo+":"+name)
		}
	}
	return strings.Join(t, ",")
}

func pullTag(tag intl.Tag) string {
	if len(tag.Names) == 0 {
		return ""
	}
	names := make([]string, 0, len(tag.Names))
	for name := range tag.Names {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		if len(names[i]) < len(names[j]) {
			return names[i] < names[j]
		}
		return false
	})
	return names[0]
}

func (m *MatrixBuilder) platforms(tag intl.Tag) string {
	ps := make([]string, 0)
	for p := range tag.Platforms {
		pl := p.String()
		if m.AllowedPlatforms == nil || len(m.AllowedPlatforms) == 0 || m.AllowedPlatforms[pl] {
			ps = append(ps, pl)
		}
	}
	return strings.Join(ps, ",")
}

func getTagsToUpdate(src []intl.Tag, dst []intl.Tag) []intl.Tag {
	aux := make(map[string]intl.Tag)
	for _, dt := range dst {
		for name := range dt.Names {
			aux[name] = dt
		}
	}
	tags := make([]intl.Tag, 0)
	for _, st := range src {
		for name := range st.Names {
			dt, exist := aux[name]
			if !exist || dt.LatestUpdate.Before(st.LatestUpdate) {
				tags = append(tags, st)
				break
			}
		}
	}
	return tags
}
