package main

import (
	"fmt"
	"math"
	"sync"

	"github.com/wilhelm-murdoch/jumpy/cmd/scraper/library"
	"github.com/wilhelm-murdoch/jumpy/cmd/scraper/logger"
	"github.com/wilhelm-murdoch/jumpy/cmd/scraper/models"
)

const (
	RootUrl = "https://wheresthejump.com/full-movie-list/"
	BaseDir = "../.."
)

func main() {
	feed, err := library.GetMovieListFromUrl(RootUrl)

	if err != nil {
		logger.Error(err.Error())
	}

	const batchSize int = 25

	shift := 0
	numUrls := len(feed.Movies)
	numBatches := int(math.Ceil(float64(numUrls)/float64(batchSize))) - 1

	logger.Info("processing %d urls in batches of %d each", numUrls, batchSize)

	for i := 0; i <= numBatches; i++ {
		logger.Info("==> Processing Batch #%d:", i)

		offset, limit := shift, shift+batchSize

		if limit > numUrls {
			limit = numUrls
		}

		currentBatch := feed.Movies[offset:limit]
		shift += batchSize
		chanError := make(chan error)
		chanFinished := make(chan bool)
		batchErrors := make([]error, 0)

		go func() {
			for {
				select {
				case err := <-chanError:
					batchErrors = append(batchErrors, err)
				case <-chanFinished:
					close(chanError)
					close(chanFinished)
					return
				}
			}
		}()

		var urlProcessingWaitGroup sync.WaitGroup
		urlProcessingWaitGroup.Add(len(currentBatch))

		for _, movie := range currentBatch {
			go func(movie models.Movie) {
				err := library.GetMovieDetailsFromUrl(&movie)
				if err != nil {
					chanError <- err
					return
				}

				feed.Add(&movie)

				err = feed.Save(fmt.Sprintf("dist/movies/%s.json", movie.Id), movie)
				if err != nil {
					logger.Error(err.Error())
				}

				logger.Info("... generated movie: dist/movies/%s.json", movie.Id)
				defer urlProcessingWaitGroup.Done()
			}(movie)
		}

		urlProcessingWaitGroup.Wait()

		chanFinished <- true
	}

	tags := feed.GetDistinctTags()
	logger.Info("processing %d distinct tags and assigning movies", len(tags))
	for _, t := range tags {
		filtered := feed.FilterMoviesByTag(&t)
		logger.Info("... generated tag: dist/tags/%s.json", t.Id)
		feed.Save(fmt.Sprintf("dist/tags/%s.json", t.Id), filtered)
	}

	feed.Save("dist/movies.json", feed.Movies)
	logger.Info("... generated movie index: dist/movies.json")

	feed.Save("dist/tags.json", tags)
	logger.Info("... generated tag index: dist/tags.json")

	logger.Info("all done; exiting ...")
}
