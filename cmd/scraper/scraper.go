package main

import (
	_ "tools/logger"
)

const RootUrl = "https://wheresthejump.com/full-movie-list/"

func main() {
	logger.Err("sup")
	// feed, err := util.GetMovieListFromUrl(RootUrl)
	// if err != nil {
	// 	logger.Err.Fatal(err)
	// }

	// const batchSize int = 25

	// shift := 0
	// numUrls := len(feed.Movies)
	// numBatches := int(math.Ceil(float64(numUrls)/float64(batchSize))) - 1

	// logger.Inf.Printf("processing %d urls in batches of %d each", numUrls, batchSize)

	// for i := 0; i <= numBatches; i++ {
	// 	logger.Inf.Printf("==> Processing Batch #%d:", i)

	// 	offset, limit := shift, shift+batchSize

	// 	if limit > numUrls {
	// 		limit = numUrls
	// 	}

	// 	currentBatch := feed.Movies[offset:limit]
	// 	shift += batchSize
	// 	chanError := make(chan error)
	// 	chanFinished := make(chan bool)
	// 	batchErrors := make([]error, 0)

	// 	go func() {
	// 		for {
	// 			select {
	// 			case err := <-chanError:
	// 				batchErrors = append(batchErrors, err)
	// 			case <-chanFinished:
	// 				close(chanError)
	// 				close(chanFinished)
	// 				return
	// 			}
	// 		}
	// 	}()

	// 	var urlProcessingWaitGroup sync.WaitGroup
	// 	urlProcessingWaitGroup.Add(len(currentBatch))

	// 	for _, movie := range currentBatch {
	// 		go func(movie movie.Movie) {
	// 			err := util.GetMovieDetailsFromUrl(movie)
	// 			if err != nil {
	// 				chanError <- err
	// 				return
	// 			}

	// 			err = movie.Save(fmt.Sprintf("json/movies/%s.json", movie.Id))
	// 			if err != nil {
	// 				logger.Err.Fatal(err)
	// 			}

	// 			logger.Inf.Printf("... generated: %s.json", movie.Id)
	// 			defer urlProcessingWaitGroup.Done()
	// 		}(movie)
	// 	}

	// 	urlProcessingWaitGroup.Wait()

	// 	chanFinished <- true
	// }

	// feed.Save("json/movies.json")

	// for _, thing := range feed.Tags {
	// 	fmt.Println(thing.Movies)
	// }

	// util.WriteJsonToFile("json/tags.json", feed.Tags)

	// tags := make([]Tag, 0)
	// for _, movie := range feed.Movies {
	// 	for _, tag := range movie.Tags {
	// 		if !ExistsByKey(tags, tag.Id) {
	// 			tags = append(tags, tag)
	// 		}
	// 	}
	// }

	// writeJsonToFile("json/tags.json", tags)
}

// func ExistsByKey(tags []tag.Tag, value string) bool {
// 	for _, tag := range tags {
// 		if tag.Id == value {
// 			return true
// 		}
// 	}
// 	return false
// }
