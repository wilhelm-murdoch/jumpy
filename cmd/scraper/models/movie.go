package models

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

type Movie struct {
	Id                      string      `json:"id"`
	Title                   string      `json:"title"`
	JumpScareCount          int         `json:"jump_scare_count,omitempty"`
	Runtime                 int         `json:"runtime,omitempty"`
	Synopsis                string      `json:"synopsis,omitempty"`
	ReleaseYear             int         `json:"release_year,omitempty"`
	CoverUrl                string      `json:"cover_url,omitempty"`
	SourceUrl               string      `json:"source_url"`
	DetailsUrl              string      `json:"details_url,omitempty"`
	Directors               []string    `json:"directors,omitempty"`
	JumpScares              []JumpScare `json:"jump_scares,omitempty"`
	Tags                    []Tag       `json:"tags,omitempty"`
	Reviews                 []Review    `json:"reviews,omitempty"`
	ContentRating           string      `json:"content_rating,omitempty"`
	JumpScareSrtUrl         string      `json:"jump_scare_srt_url,omitempty"`
	JumpScareSpoilersSrtUrl string      `json:"jump_scare_spoilers_srt_url,omitempty"`
}

func NewMovie(title string, release int, url string) *Movie {
	id := slug.Make(fmt.Sprintf("%s-%d", title, release))
	return &Movie{
		Id:          id,
		Title:       title,
		ReleaseYear: release,
		SourceUrl:   url,
		DetailsUrl:  fmt.Sprintf("https://jumpy.wilhelm.codes/movies/%s.json", id),
	}
}

func (m *Movie) SaveSrt(path string, spoilers bool) error {
	var output []string
	var prefix string = "Minor "

	subtitle := "jump scare ahead!"
	for i, j := range m.JumpScares {
		if spoilers {
			subtitle = "- " + j.Spoiler
		}

		if j.Major {
			prefix = "Major "
		}

		output = append(
			output,
			fmt.Sprint(i+1),
			j.TimeStart+" --> "+j.TimeStop,
			prefix+subtitle+"\n",
		)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.WriteString(strings.Join(output, "\n"))
	if err != nil {
		return err
	}

	defer file.Close()

	m.JumpScareSrtUrl = fmt.Sprintf("https://jumpy.wilhelm.codes/movies/%s.srt", m.Id)
	m.JumpScareSpoilersSrtUrl = fmt.Sprintf("https://jumpy.wilhelm.codes/movies/%s-spoilers.srt", m.Id)

	return nil
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
	id := slug.Make(name)

	m.Tags = append(m.Tags, Tag{
		Id:         id,
		Name:       strings.Trim(name, " "),
		DetailsUrl: fmt.Sprintf("https://jumpy.wilhelm.codes/tags/%s.json", id),
	})
}

func (m *Movie) AddReview(name string, url string) {
	m.Reviews = append(m.Reviews, Review{
		Name: name,
		Url:  url,
	})
}

func (m *Movie) HasTag(tag *Tag) bool {
	for _, t := range m.Tags {
		if t.Id == tag.Id {
			return true
		}
	}
	return false
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
