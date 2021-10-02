package main

import (
	"encoding/json"
	"flag"
	"fmt"
	srv "github.com/s6on/mirror-docker-tags-action/internal/service"
	"log"
	"os"
	"strings"
)

func main() {
	var from string
	var to string
	var extraRegistry string
	var allowedPs string
	var updateAll bool

	flag.StringVar(&from, "from", "", "Comma separate docker repositories to mirror the tags from")
	flag.StringVar(&to, "to", os.Getenv("GITHUB_REPOSITORY_OWNER"), "Docker repositories to mirror the tags into")
	flag.StringVar(&extraRegistry, "extraRegistry", "", "Extra registry to mirror the tags into")
	flag.StringVar(&allowedPs, "allowedPlatforms", "", "Comma separate list of allowed Platforms")
	flag.BoolVar(&updateAll, "updateAll", false, "Update all tags")

	flag.Parse()

	if len(from) == 0 || len(to) == 0 {
		log.Fatalln("Missing parameter")
	}

	m := srv.MatrixBuilder{
		From:             getRepos(from),
		To:               to,
		ExtraRegistry:    extraRegistry,
		UpdateAll:        updateAll,
		AllowedPlatforms: allowedPlatforms(allowedPs),
	}

	matrix := m.Get()
	if len(matrix.Include) == 0 {
		return
	}
	fm, err := json.Marshal(matrix)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("::set-output name=matrix::%s", fm)
}

func getRepos(from string) map[string][]string {
	rt := make(map[string][]string)
	repos := strings.Split(from, "]")
	for _, repo := range repos {
		if strings.Contains(repo, "[") {
			r, t := getRepo(repo)
			rt[r] = t
		}
	}
	return rt
}

func getRepo(rt string) (string, []string) {
	split := strings.Split(rt, "[")
	if len(split) != 2 {
		log.Fatal("Invalid input")
	}
	repo := strings.ReplaceAll(split[0], ",", "")
	return repo, strings.Split(split[1], ",")
}

func allowedPlatforms(platformss string) map[string]bool {
	platforms := strings.Split(platformss, ",")
	if len(platforms) == 0 {
		return nil
	}
	ap := make(map[string]bool)
	for _, platform := range platforms {
		ap[platform] = true
	}
	return ap
}
