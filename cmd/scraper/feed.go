package main

import (
	"github.com/wilhelm-murdoch/jumpscare-api/movie"
	"github.com/wilhelm-murdoch/jumpscare-api/tag"
	"github.com/wilhelm-murdoch/jumpscare-api/util"
)

type Feed struct {
	Movies []movie.Movie `json:"movies"`
	Tags   []tag.Tag     `json:"tags"`
}

func NewFeed() *Feed {
	return &Feed{}
}

func (f *Feed) Add(movie *movie.Movie) {
	f.Movies = append(f.Movies, *movie)
}

func (f *Feed) Save(path string) error {
	return util.WriteJsonToFile(path, f.Movies)
}
