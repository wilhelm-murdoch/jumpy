package library

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wilhelm-murdoch/go-batch"
	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/logger"
	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/models"

	"github.com/PuerkitoBio/goquery"
)

func GetDocumentFromUrl(url string) (*goquery.Document, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	var backoff time.Duration

	maxAttempts := 10
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		response, err := client.Get(url)
		if err != nil {
			break
		}

		switch response.StatusCode {
		case 429:
			if attempt >= maxAttempts {
				logger.Warning("could not get url %s after %d attempts; skipping ...", url, attempt)
				return goquery.NewDocumentFromReader(response.Body)
			}

			backoff = time.Duration(attempt) * time.Second
			logger.Warning("got rate-limited on url %s; waiting another %d seconds", url, (backoff / time.Second))
			time.Sleep(backoff)
		case 200:
			return goquery.NewDocumentFromReader(response.Body)
		case 404:
			logger.Warning("url %s could not be found; skipping ...", url)
			return goquery.NewDocumentFromReader(response.Body)
		}

		defer response.Body.Close()
	}

	return nil, errors.New("this should not happen, but here we are")
}

func GetMovieListFromUrl(url string) (batch.Iterator[models.Movie], error) {
	iterator := batch.Iterator[models.Movie]{}

	document, err := GetDocumentFromUrl(url)
	if err != nil {
		return iterator, err
	}

	document.Find("div.entry-content > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		node := s.Find("td:nth-child(1) a")

		url, ok := node.Attr("href")
		if ok {
			title := node.Text()
			release, _ := strconv.Atoi(s.Find("td:nth-child(3)").Text())
			movie := models.NewMovie(title, release, url)

			iterator.Push(&movie)
		}
	})

	return iterator, nil
}

func GetMovieDetailsFromUrl(movie models.Movie) (models.Movie, error) {
	document, err := GetDocumentFromUrl(movie.SourceUrl)
	if err != nil {
		return movie, err
	}

	document.Find("div.inside-article").Each(func(i int, s *goquery.Selection) {
		s.Find("div.entry-content > h3 ~ p").Each(func(i int, s *goquery.Selection) {
			pattern := regexp.MustCompile(`^(\d{2}:\d{2}:\d{2}) – (.*).$`)
			matches := pattern.FindStringSubmatch(s.Text())

			if len(matches) == 3 {
				var major bool = false
				if s.Has("strong").Length() >= 1 {
					major = true
				}

				movie.AddJumpScare(matches[1], matches[2], major)
			}
		})

		s.Find("div.video-info-grid-column > p").Each(func(i int, s *goquery.Selection) {
			switch strings.Fields(s.Text())[0] {
			case "Reviews:":
				s.Find("a").Each(func(i int, s *goquery.Selection) {
					href, ok := s.Attr("href")
					if ok {
						movie.AddReview(s.Text(), href)
					}
				})
			case "Director:":
				for _, d := range strings.Split(strings.Replace(s.Text(), "Director:", "", 1), ",") {
					movie.AddDirector(d)
				}
			case "Rating:":
				movie.SetContentRating(s.Text())
			case "Runtime:":
				movie.SetRuntimeFromPattern(regexp.MustCompile(`^Runtime:\s+(\d+)\sminutes$`), s.Text())
			case "Tags:":
				s.Find("a").Each(func(i int, s *goquery.Selection) {
					movie.AddTag(s.Text())
				})
			}
		})
	})

	return movie, nil
}
