package dockerhub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Active bool

func (s *Active) UnmarshalJSON(data []byte) error {
	*s = string(data) == "\"active\""
	return nil
}

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) <= 4 {
		return nil
	}
	vt, err := time.Parse("2006-01-02T15:04:05.00", string(data)[1:23])
	if err == nil {
		*t = Time(vt)
	}
	return err
}

type TagsResponse struct {
	Count uint
	Next  string
	Tags  []Tag `json:"results"`
}

type Tag struct {
	Active      Active `json:"tag_status"`
	LastUpdated Time   `json:"last_updated"`
	Name        string
	Images      []Image
}

type Image struct {
	Active       Active `json:"status"`
	Digest       string
	Os           string
	Architecture string
	Variant      string
}

func (t *Tag) Id() string {
	uniq := make(map[string]bool)
	digs := make([]string, 0)
	for _, img := range t.Images {
		if uniq[img.Digest] {
			continue
		}
		uniq[img.Digest] = true
		digs = append(digs, img.Digest)
	}
	sort.Strings(digs)
	return strings.Join(digs, ",")
}

func Get(url string) TagsResponse {
	rs, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer rs.Body.Close()
	tagsRs := TagsResponse{}
	err = json.NewDecoder(rs.Body).Decode(&tagsRs)
	if err != nil {
		log.Fatalln(err)
	}
	return tagsRs
}

func GetTags(repo string, pageSize uint) TagsResponse {
	return Get(fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags?page_size=%d", cleanup(repo), pageSize))
}

func cleanup(repo string) string {
	if strings.HasPrefix(repo, "_/") {
		repo = strings.Replace(repo, "_/", "library/", 1)
	} else if !strings.Contains(repo, "/") {
		repo = "library/" + repo
	}
	return repo
}
