package main

import (
	"fmt"
	"math"
	"sync"

	"github.com/wilhelm-murdoch/jumpscare-api/cmd/scraper/library"
	"github.com/wilhelm-murdoch/jumpscare-api/cmd/scraper/logger"
	"github.com/wilhelm-murdoch/jumpscare-api/cmd/scraper/models"
)

const RootUrl = "https://wheresthejump.com/full-movie-list/"

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

				err = feed.Save(fmt.Sprintf("../../build/json/movies/%s.json", movie.Id), movie)
				if err != nil {
					logger.Error(err.Error())
				}

				logger.Info("... generated: %s.json", movie.Id)
				defer urlProcessingWaitGroup.Done()
			}(movie)
		}

		urlProcessingWaitGroup.Wait()

		chanFinished <- true
	}

	feed.Save("json/movies.json", feed.Movies)

	tags := feed.GetDistinctTags()
	for _, t := range tags {
		filtered := feed.FilterMoviesByTag(&t)
		feed.Save(fmt.Sprintf("../../build/json/tags/%s.json", t.Id), filtered)
	}

	feed.Save("../../build/json/tags.json", tags)
}
