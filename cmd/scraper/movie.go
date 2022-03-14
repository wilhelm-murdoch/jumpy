package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wilhelm-murdoch/jumpscare-api/tag"
	"github.com/wilhelm-murdoch/jumpscare-api/util"

	"github.com/gosimple/slug"
)

type Movie struct {
	Id             string      `json:"id"`
	Title          string      `json:"title"`
	JumpScareCount int         `json:"jump_scare_count,omitempty"`
	Runtime        int         `json:"runtime,omitempty"`
	Synopsis       string      `json:"synopsis,omitempty"`
	ReleaseYear    int         `json:"release_year,omitempty"`
	CoverUrl       string      `json:"cover_url,omitempty"`
	SourceUrl      string      `json:"source_url"`
	DetailsUrl     string      `json:"details_url,omitempty"`
	Directors      []string    `json:"directors,omitempty"`
	JumpScares     []JumpScare `json:"jump_scares,omitempty"`
	Tags           []tag.Tag   `json:"tags,omitempty"`
	Reviews        []Review    `json:"reviews,omitempty"`
	ContentRating  string      `json:"content_rating,omitempty"`
}

func NewMovie(title string, release int, url string) *Movie {
	return &Movie{
		Id:          slug.Make(fmt.Sprintf("%s-%d", title, release)),
		Title:       title,
		ReleaseYear: release,
		SourceUrl:   url,
	}
}

func (m *Movie) Save(path string) error {
	return util.WriteJsonToFile(path, m)
}

func (m *Movie) SaveSrt(path string) error {
	var output string
	for _, jumpscare := range m.JumpScares {
		fmt.Println(jumpscare)
	}
	return util.WriteJsonToFile(path, output)
}

func (m *Movie) AddDirector(director string) {
	m.Directors = append(m.Directors, strings.Trim(director, " "))
}

func (m *Movie) AddJumpScare(timestamp string, spoiler string, major bool) {
	parsed, _ := time.Parse("03:04:05", timestamp)
	rewind, _ := time.ParseDuration("5s")
	timeStart := parsed.Add(-rewind)

	m.JumpScareCount++

	m.JumpScares = append(m.JumpScares, JumpScare{
		TimeStart: fmt.Sprintf("%02d:%02d:%02d", timeStart.Hour(), timeStart.Minute(), timeStart.Second()),
		TimeStop:  timestamp,
		Spoiler:   spoiler,
		Major:     major,
	})
}

func (m *Movie) AddTag(name string) {
	m.Tags = append(m.Tags, tag.Tag{
		Id:   slug.Make(name),
		Name: strings.Trim(name, " "),
	})
}

func (m *Movie) AddReview(name string, url string) {
	m.Reviews = append(m.Reviews, Review{
		Name: name,
		Url:  url,
	})
}

func (m *Movie) SetContentRating(rating string) {
	m.ContentRating = strings.Trim(strings.Replace(rating, "Rating: ", "", 1), " ")
}

func (m *Movie) SetRuntimeFromPattern(pattern *regexp.Regexp, text string) {
	matches := pattern.FindStringSubmatch(text)

	if len(matches) == 2 {
		m.Runtime, _ = strconv.Atoi(matches[1])
	}
}

type JumpScare struct {
	TimeStart string `json:"time_start"`
	TimeStop  string `json:"time_stop"`
	Spoiler   string `json:"spoiler"`
	Major     bool   `json:"major"`
}

type Review struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
