package main

type Tag struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Movies []Movie `json:"movies,omitempty"`
}

func (t *Tag) AddMovie(movie *Movie) bool {
	var found bool = false
	for _, m := range t.Movies {
		if m.Id == movie.Id {
			found = true
		}
	}

	if !found {
		t.Movies = append(t.Movies, *movie)
	}

	return found
}
