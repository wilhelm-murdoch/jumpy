package models

import (
	"encoding/json"
	"os"
)

type Feed struct {
	Movies []Movie `json:"movies"`
}

func NewFeed() *Feed {
	return &Feed{}
}

func (f *Feed) FilterMoviesByTag(tag *Tag) *Feed {
	feed := NewFeed()

	for _, m := range f.Movies {
		if m.HasTag(tag) {
			feed.Push(m)
		}
	}

	return feed
}

func (f *Feed) GetDistinctTags() []Tag {
	distinct := make(map[string]Tag, 0)

	for _, m := range f.Movies {
		for _, t := range m.Tags {
			if _, ok := distinct[t.Id]; !ok {
				distinct[t.Id] = t
			}
		}
	}

	tags := make([]Tag, 0, len(distinct))
	for _, t := range distinct {
		tags = append(tags, t)
	}

	return tags
}

func (f *Feed) Push(movie Movie) {
	f.Movies = append(f.Movies, movie)
}

func (f *Feed) Length() int {
	return len(f.Movies)
}

func (f *Feed) Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)

	err = encoder.Encode(object)
	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}
