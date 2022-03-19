package main

// job retries; with limits, track retry attempts
// return final result from processor.Execute
// implement push pop each
// update models to use generics where relevant
// create "saver" interface with Feed | Movie types?
// error types; retry, failed, etc ...
// return completed jobs as a collection .each .push. .pop .map

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/processor"
)

// const (
// 	RootUrl = "https://wheresthejump.com/full-movie-list/"
// 	BaseDir = "../.."
// )

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var things []string
	for i := 0; i < 10; i++ {
		things = append(things, randSeq(15))
	}

	ps := processor.NewProcessor(things)
	results := ps.Execute(func(job *processor.Job[string]) (*processor.Job[string], error) {
		fmt.Println("Processing:     ", job.Body)
		job.Body += " edit"
		return job, nil
	})

	fmt.Println()
	fmt.Println("Finished:")
	res := processor.Results[string]{}

	res.Concat(results).Each(func(i int, j *processor.Job[string]) {
		fmt.Println(j.Body)
	})
}

// func main() {
// 	summary, err := library.GetMovieListFromUrl(RootUrl)
// 	if err != nil {
// 		logger.Error(err.Error())
// 	}

// 	ps := processor.NewProcessor(summary.Movies, 100)
// 	ps.Execute(handlers.HandleMovieDetails)
// }

// movie.SaveSrt("dist"+fmt.Sprintf("/downloads/srt/%s-spoilers.srt", movie.Id), true)
// logger.Info("... generated movie: dist/movies/%s.json", movie.Id)

// movie.SaveSrt("dist"+fmt.Sprintf("/downloads/srt/%s.srt", movie.Id), false)
// logger.Info("... generated movie: dist/movies/%s.json", movie.Id)
// tags := feed.GetDistinctTags()
// logger.Info("processing %d distinct tags and assigning movies", len(tags))
// for _, t := range tags {
// 	filtered := feed.FilterMoviesByTag(&t)
// 	logger.Info("... generated tag: dist/tags/%s.json", t.Id)
// 	feed.Save(fmt.Sprintf("dist/tags/%s.json", t.Id), filtered)
// }

// summary.Save("dist/movies.json", summary.Movies)
// logger.Info("... generated movie index: dist/movies.json")

// feed.Save("dist/tags.json", tags)
// logger.Info("... generated tag index: dist/tags.json")

// logger.Info("all done; exiting ...")
