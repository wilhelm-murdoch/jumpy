package main

// job retries; with limits, track retry attempts
// create "saver" interface with Feed | Movie types?
// error types; retry, failed, etc ...
// create signal struct supporting various signals
// convert Result struct to Result type

import (
	"os"

	"github.com/wilhelm-murdoch/go-batch"
	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/handlers"
	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/library"
	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/logger"
)

func main() {
	// Let's hit https://wheresthejump.com/ to get a list of all movies. This won't
	// have all the information we need, but we can prime each record with some
	// basic information.
	summary, err := library.GetMovieListFromUrl("https://wheresthejump.com/full-movie-list/")
	if err != nil {
		logger.Error(err.Error())
	}

	os.Exit(0)

	// Here we have a signal handler so we can listen to anything interesting as
	// we process all the movies in batches.
	sh := batch.SignalHandler(func(e error) {
		logger.Warning(e.Error())
	})

	// Create a batch processor to scrape the target movie pages for more detailed
	// information and return the results.
	ps := batch.NewProcessor(summary.Data, sh)
	ps.Execute(handlers.HandleMovieDetails)

	// This is a bit janky at the moment. I'd like to just pass in the previous
	// result set and pull the trigger, but unfortunately I'm being a bit lazy.
	// So, I just rebuild a new slice of models.Movie sourced from the batch.Job
	// set from the results.
	// movies := batch.Iterator[models.Movie]{}
	// results.Each(func(i int, j *batch.Job[models.Movie]) {
	// 	logger.Info("... generated movie: dist/movies/%s.json", j.Body.Id)
	// 	j.Body.Save(fmt.Sprintf("dist/movies/%s.json", j.Body.Id), j.Body)

	// 	movies.Push(&j.Body)
	// })

	// Trigger a new batch processor to try to ingest a movie poster for each
	// movie from https://www.movieposterdb.com/. Unfortunately, they go pretty
	// heavy on the rate limiting, so each batch is only 5 jobs.
	// ps = batch.NewProcessor(movies.Data, batch.BatchSize(5), sh)
	// results = ps.Execute(handlers.HandleMoviePoster)

	// full := batch.Iterator[models.Movie]{}
	// results.Each(func(i int, j *batch.Job[models.Movie]) {
	// 	full.Push(j.Body)
	// })

	// tags := full.GetDistinctTags()
	// logger.Info("processing %d distinct tags and assigning movies", len(tags))
	// for _, t := range tags {
	// 	filtered := full.FilterMoviesByTag(&t)
	// 	logger.Info("... generated tag: dist/tags/%s.json", t.Id)
	// 	full.Save(fmt.Sprintf("dist/tags/%s.json", t.Id), filtered)
	// }

	// summary.Save("dist/movies.json", summary.Data)
	// logger.Info("... generated movie index: dist/movies.json")

	// full.Save("dist/tags.json", tags)
	// logger.Info("... generated tag index: dist/tags.json")

	// logger.Info("all done; exiting ...")
}
