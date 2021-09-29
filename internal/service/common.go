package service

import (
	intl "github.com/s6on/mirror-docker-tags-action/internal"
	dckr "github.com/s6on/mirror-docker-tags-action/internal/dockerhub"
	"sort"
	"strings"
)

func Matrix(from map[string][]string, to string, extraRegistry string, updateAll bool) intl.Matrix {
	parms := make([]intl.Params, 0)
	for repo, tags := range from {
		ot := dckr.Get(repo, tags)
		trn := toRepoName(to, repo)
		var nt []intl.Tag
		if updateAll {
			nt = ot
		} else {
			tt := dckr.Get(trn, tags)
			nt = getTagsToUpdate(ot, tt)
		}
		for _, tag := range nt {
			parms = append(parms, param(tag, repo, trn, extraRegistry))
		}
	}
	return intl.Matrix{Include: parms}
}

func toRepoName(to string, repo string) string {
	return to + "/" + strings.ReplaceAll(repo, "/", "-")
}

func param(tag intl.Tag, fromRepo string, toRepo string, extra string) intl.Params {
	return intl.Params{
		BaseImg:   fromRepo + ":" + pullTag(tag),
		Tags:      tags(tag, toRepo, extra),
		Platforms: platforms(tag)}
}

func tags(tag intl.Tag, toRepo string, extra string) string {
	t := make([]string, 0)
	for name := range tag.Names {
		t = append(t, toRepo+":"+name)
		if len(extra) > 0 {
			t = append(t, extra+"/"+toRepo+":"+name)
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

func platforms(tag intl.Tag) string {
	ps := make([]string, 0)
	for p := range tag.Platforms {
		ps = append(ps, p.String())
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
